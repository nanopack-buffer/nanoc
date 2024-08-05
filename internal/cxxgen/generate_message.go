package cxxgen

import (
	"errors"
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"nanoc/internal/pathutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/iancoleman/strcase"
)

// Options are parameters that can be tweaked to alter codegen.
type Options struct {
	BaseDirectoryPath   string
	OutputDirectoryPath string

	FormatterPath string
	FormatterArgs []string

	// The absolute path to the directory where the factory file should be put in
	// This is an empty string when it is not requested.
	MessageFactoryPath string

	// The namespace to which put the generated messages should be put under.
	// For nested namespaces, the outermost namespace will appear earlier in the slice.
	// For example, My.Namespace will be stored as ["My", "Namespace"] in the slice.
	Namespaces []string
}

func GenerateMessageClass(msgSchema *npschema.Message, opts Options) error {
	ng := numberGenerator{}
	gm := generator.MessageCodeGeneratorMap{
		datatype.Int8:    ng,
		datatype.Int32:   ng,
		datatype.Int64:   ng,
		datatype.UInt8:   ng,
		datatype.UInt32:  ng,
		datatype.UInt64:  ng,
		datatype.Double:  ng,
		datatype.String:  stringGenerator{},
		datatype.Bool:    boolGenerator{},
		datatype.Message: messageGenerator{},
		datatype.Any:     anyGenerator{},
	}
	gm[datatype.Enum] = enumGenerator{gm}
	gm[datatype.Optional] = optionalGenerator{gm}
	gm[datatype.Array] = arrayGenerator{gm}
	gm[datatype.Map] = mapGenerator{gm}

	var wg sync.WaitGroup
	var errs []error
	if msgSchema.IsInherited {
		errs = make([]error, 4)
	} else {
		errs = make([]error, 2)
	}

	wg.Add(1)
	go func() {
		errs[0] = generateMessageHeaderFile(msgSchema, gm, opts)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errs[1] = generateMessageImplFile(msgSchema, gm, opts)
		wg.Done()
	}()

	if msgSchema.IsInherited {
		wg.Add(1)
		go func() {
			errs[2] = generateChildMessageFactoryHeaderFile(msgSchema, opts)
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			errs[3] = generateChildMessageFactoryImplFile(msgSchema, opts)
			wg.Done()
		}()
	}

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return errors.Join(errs...)
		}
	}

	return nil
}

func GenerateMessageFactory(schemas []*npschema.Message, opts Options) error {
	var wg sync.WaitGroup
	errs := make([]error, 2)

	wg.Add(1)
	go func() {
		errs[0] = generateMessageFactoryHeaderFile(opts)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errs[1] = generateMessageFactoryImplFile(schemas, opts)
		wg.Done()
	}()

	wg.Wait()

	if errs[0] != nil || errs[1] != nil {
		return errors.Join(errs...)
	}

	return nil
}

func findGeneratorForField(f npschema.MessageField, gm generator.MessageCodeGeneratorMap) generator.DataTypeMessageCodeGenerator {
	if f.IsSelfReferencing() {
		return gm[f.Type.ElemType.Kind]
	}
	return gm[f.Type.Kind]
}

func generateMessageHeaderFile(msgSchema *npschema.Message, gm generator.MessageCodeGeneratorMap, opts Options) error {
	info := messageHeaderFileTemplateInfo{
		Namespace:        strings.Join(opts.Namespaces, cxxSymbolMemberOf),
		MessageName:      msgSchema.Name,
		TypeID:           msgSchema.TypeID,
		HasParentMessage: msgSchema.HasParentMessage,
		IncludeGuardName: strcase.ToScreamingSnake(msgSchema.Name) + "_NP_HXX",
		IsInherited:      msgSchema.IsInherited,
	}

	libimp := map[string]struct{}{}
	for _, field := range msgSchema.DeclaredFields {
		if field.IsSelfReferencing() {
			libimp["memory"] = struct{}{}
			continue
		}

		var collectImport func(t datatype.DataType)
		collectImport = func(t datatype.DataType) {
			switch t.Kind {
			case datatype.String:
				libimp["string"] = struct{}{}

			case datatype.Map:
				libimp["unordered_map"] = struct{}{}
				collectImport(*t.ElemType)

			case datatype.Optional:
				libimp["optional"] = struct{}{}
				collectImport(*t.ElemType)

			case datatype.Any:
				libimp["nanopack/any.hxx"] = struct{}{}

			case datatype.Message:
				if field.Type.Schema.(*npschema.Message).IsInherited {
					libimp["memory"] = struct{}{}
				}

			case datatype.Array:
				libimp["vector"] = struct{}{}

			default:
				break
			}
		}

		collectImport(field.Type)
	}

	if msgSchema.HasParentMessage {
		info.ParentMessageName = msgSchema.ParentMessage.Name
	}

	for k := range libimp {
		info.LibraryImports = append(info.LibraryImports, k)
	}
	for _, t := range msgSchema.ImportedTypes {
		p, err := resolveSchemaImportPath(t, msgSchema)
		if err != nil {
			return err
		}
		info.RelativeImports = append(info.RelativeImports, p)

		// if this type is inherited (polymorphic) and is imported because it is used by one of the fields
		// then its factory needs to be imported as well,
		// because it will be used when reading fields that use this polymorphic type to instantiate the correct type.
		if ms, ok := t.(*npschema.Message); ok && ms.IsInherited {
			header := fmt.Sprintf("make_%v%v", strcase.ToSnake(ms.Name), extHeaderFile)
			info.RelativeImports = append(info.RelativeImports, strings.Replace(p, filepath.Base(p), header, 1))
		}
	}

	for _, f := range msgSchema.InheritedFields {
		g := findGeneratorForField(f, gm)
		info.ConstructorParameters = append(info.ConstructorParameters, g.ConstructorFieldParameter(f))
	}

	for _, f := range msgSchema.DeclaredFields {
		g := findGeneratorForField(f, gm)
		info.FieldDeclarationLines = append(info.FieldDeclarationLines, g.FieldDeclaration(f))
		info.ConstructorParameters = append(info.ConstructorParameters, g.ConstructorFieldParameter(f))

		if f.Type.Kind == datatype.Message && f.Type.Schema.(*npschema.Message).IsInherited {
			// this field stores a polymorphic type which requires a unique_ptr to hold the value
			// a getter is needed to expose the value as a reference
			l := fmt.Sprintf("[[nodiscard]] %v &get_%v() const;", f.Type.Identifier, strcase.ToSnake(f.Name))
			info.FieldGetters = append(info.FieldGetters, l)
		}
	}

	tmpl, err := template.New(templateNameMessageHeaderFile).
		Funcs(template.FuncMap{
			"join": strings.Join,
		}).
		Parse(messageHeaderFile)

	if err != nil {
		return err
	}

	fname := filepath.Base(msgSchema.SchemaPath)
	fname = strcase.ToSnake(strings.TrimSuffix(fname, filepath.Ext(fname))) + extHeaderFile
	outPath := pathutil.ResolveCodeOutputPathForSchema(msgSchema, opts.BaseDirectoryPath, opts.OutputDirectoryPath, fname)

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, info)
	if err != nil {
		return err
	}

	err = formatCode(outPath, opts.FormatterPath, opts.FormatterArgs...)
	if err != nil {
		return err
	}

	return nil
}

func generateMessageImplFile(msgSchema *npschema.Message, gm generator.MessageCodeGeneratorMap, opts Options) error {
	fname := filepath.Base(msgSchema.SchemaPath)
	fname = strcase.ToSnake(strings.TrimSuffix(fname, filepath.Ext(fname))) + extImplFile

	// the message header byte size includes 4 bytes for the type ID and 4 bytes to store the byte size of each field
	npHeaderByteSize := (len(msgSchema.AllFields) + 1) * 4

	info := messageImplFileTemplateInfo{
		Namespace:               strings.Join(opts.Namespaces, cxxSymbolMemberOf),
		HeaderName:              strings.Replace(fname, extImplFile, extHeaderFile, 1),
		MessageName:             msgSchema.Name,
		HasParentMessage:        msgSchema.HasParentMessage,
		ReadPtrStart:            npHeaderByteSize,
		ConstructorParameters:   nil,
		SuperConstructorArgs:    nil,
		FieldInitializers:       nil,
		HeaderSize:              npHeaderByteSize,
		FieldReadCodeFragments:  nil,
		FieldWriteCodeFragments: nil,
	}
	if msgSchema.HasParentMessage {
		info.ParentMessageName = msgSchema.ParentMessage.Name
	}

	for _, f := range msgSchema.InheritedFields {
		g := findGeneratorForField(f, gm)
		info.ConstructorParameters = append(info.ConstructorParameters, g.ConstructorFieldParameter(f))
		info.SuperConstructorArgs = append(info.SuperConstructorArgs, strcase.ToSnake(f.Name))
	}

	for _, f := range msgSchema.DeclaredFields {
		g := findGeneratorForField(f, gm)
		info.ConstructorParameters = append(info.ConstructorParameters, g.ConstructorFieldParameter(f))
		info.FieldInitializers = append(info.FieldInitializers, g.FieldInitializer(f))

		if f.Type.Kind == datatype.Message && f.Type.Schema.(*npschema.Message).IsInherited {
			// this field stores a polymorphic type which requires a unique_ptr to hold the value
			// a getter is needed to expose the value as a reference
			s := strcase.ToSnake(f.Name)

			var l0 string
			if info.Namespace == "" {
				l0 = fmt.Sprintf("%v &%v::get_%v() const {", f.Type.Identifier, info.MessageName, s)
			} else {
				l0 = fmt.Sprintf("%v::%v &%v::%v::get_%v() const {", info.Namespace, f.Type.Identifier, info.Namespace, info.MessageName, s)
			}

			l := generator.Lines(
				l0,
				fmt.Sprintf("    return *%v;", s),
				"}")

			info.FieldGetters = append(info.FieldGetters, l)
		}
	}

	ctx := generator.NewCodeContext()
	for _, f := range msgSchema.AllFields {
		g := findGeneratorForField(f, gm)
		info.FieldReadCodeFragments = append(info.FieldReadCodeFragments, g.ReadFieldFromBuffer(f, ctx))
	}

	// new context in new block of code
	ctx = generator.NewCodeContext()
	for _, f := range msgSchema.AllFields {
		g := findGeneratorForField(f, gm)
		info.FieldWriteCodeFragments = append(info.FieldWriteCodeFragments, g.WriteFieldToBuffer(f, ctx))
	}

	tmpl, err := template.New(templateNameMessageImplFile).
		Funcs(template.FuncMap{
			"join": strings.Join,
		}).
		Parse(messageImplFile)

	if err != nil {
		return err
	}

	outPath := pathutil.ResolveCodeOutputPathForSchema(msgSchema, opts.BaseDirectoryPath, opts.OutputDirectoryPath, fname)
	fmt.Println(opts.OutputDirectoryPath)
	fmt.Println(outPath)

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, info)
	if err != nil {
		return err
	}

	err = formatCode(outPath, opts.FormatterPath, opts.FormatterArgs...)
	if err != nil {
		return err
	}

	return nil
}

func generateChildMessageFactoryHeaderFile(msgSchema *npschema.Message, opts Options) error {
	h := filepath.Base(msgSchema.SchemaPath)
	h = strcase.ToSnake(strings.TrimSuffix(h, filepath.Ext(h))) + extHeaderFile

	info := childMessageFactoryHeaderFileTemplateInfo{
		Namespace:           strings.Join(opts.Namespaces, cxxSymbolMemberOf),
		IncludeGuardName:    fmt.Sprintf("%v_FACTORY_NP_HXX", strcase.ToScreamingSnake(msgSchema.Name)),
		MessageName:         msgSchema.Name,
		MessageHeaderName:   h,
		FactoryFunctionName: fmt.Sprintf("make_%v", strcase.ToSnake(msgSchema.Name)),
	}

	tmpl, err := template.New(templateNameChildMessageFactoryHeaderFile).Parse(childMessageFactoryHeaderFile)
	if err != nil {
		return err
	}

	fname := fmt.Sprintf("make_%v%v", strcase.ToSnake(msgSchema.Name), extHeaderFile)
	outPath := pathutil.ResolveCodeOutputPathForSchema(msgSchema, opts.BaseDirectoryPath, opts.OutputDirectoryPath, fname)
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, info)
	if err != nil {
		return err
	}

	err = formatCode(outPath, opts.FormatterPath, opts.FormatterArgs...)
	if err != nil {
		return err
	}

	return nil
}

func generateChildMessageFactoryImplFile(msgSchema *npschema.Message, opts Options) error {
	fname := filepath.Base(msgSchema.SchemaPath)
	fname = fmt.Sprintf("make_%v%v", strcase.ToSnake(msgSchema.Name), extImplFile)

	info := childMessageFactoryImplFileTemplateInfo{
		Namespace:           strings.Join(opts.Namespaces, cxxSymbolMemberOf),
		Schema:              msgSchema,
		HeaderName:          strings.Replace(fname, extImplFile, extHeaderFile, 1),
		FactoryFunctionName: fmt.Sprintf("make_%v", strcase.ToSnake(msgSchema.Name)),
	}

	for _, m := range msgSchema.ChildMessages {
		p, err := resolveSchemaImportPath(m, msgSchema)
		if err != nil {
			return err
		}
		info.ChildMessageImportPaths = append(info.ChildMessageImportPaths, p)
	}

	tmpl, err := template.New(templateNameChildMessageFactoryImplFile).Parse(childMessageFactoryImplFile)
	if err != nil {
		return err
	}

	outPath := pathutil.ResolveCodeOutputPathForSchema(msgSchema, opts.BaseDirectoryPath, opts.OutputDirectoryPath, fname)
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, info)
	if err != nil {
		return err
	}

	err = formatCode(outPath, opts.FormatterPath, opts.FormatterArgs...)
	if err != nil {
		return err
	}

	return nil
}

func generateMessageFactoryHeaderFile(opts Options) error {
	info := messageFactoryHeaderFileTemplateInfo{
		Namespace: strings.Join(opts.Namespaces, cxxSymbolMemberOf),
	}

	outPath := filepath.Join(opts.MessageFactoryPath, fileNameMessageFactory+extHeaderFile)

	tmpl, err := template.New(templateNameMessageFactoryHeaderFile).Parse(messageFactoryHeaderFile)
	if err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, info)
	if err != nil {
		return err
	}

	return formatCode(outPath, opts.FormatterPath, opts.FormatterArgs...)
}

func generateMessageFactoryImplFile(schemas []*npschema.Message, opts Options) error {
	outPath := filepath.Join(opts.MessageFactoryPath, fileNameMessageFactory+extImplFile)

	info := messageFactoryImplFileTemplateInfo{
		Namespace:          strings.Join(opts.Namespaces, cxxSymbolMemberOf),
		MessageImportPaths: nil,
		MessageSchemas:     nil,
	}

	for _, s := range schemas {
		ip, err := filepath.Rel(filepath.Dir(outPath), s.SchemaPath)
		if err != nil {
			return err
		}
		ip = strings.Replace(ip, path.Ext(ip), extHeaderFile, 1)
		info.MessageImportPaths = append(info.MessageImportPaths, ip)
	}

	info.MessageSchemas = schemas

	tmpl, err := template.New(templateNameMessageFactoryImplFile).Parse(messageFactoryImplFile)
	if err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, info)
	if err != nil {
		return err
	}

	err = formatCode(outPath, opts.FormatterPath, opts.FormatterArgs...)
	if err != nil {
		return err
	}

	return nil
}

func resolveSchemaImportPath(toSchema datatype.Schema, fromSchema datatype.Schema) (string, error) {
	p, err := filepath.Rel(filepath.Dir(fromSchema.SchemaPathAbsolute()), toSchema.SchemaPathAbsolute())
	if err != nil {
		return "", err
	}
	return strings.Replace(p, path.Ext(p), extHeaderFile, 1), nil
}

func formatCode(path string, formatter string, args ...string) error {
	args = append(args, path)
	cmd := exec.Command(formatter, args...)
	return cmd.Run()
}
