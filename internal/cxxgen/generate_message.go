package cxxgen

import (
	"errors"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

// Options are parameters that can be tweaked to alter codegen.
type Options struct {
	FormatterPath string
	FormatterArgs []string

	// The absolute path to the directory where the factory file should be put in
	// This is an empty string when it is not requested.
	MessageFactoryPath string
}

func GenerateMessageClass(msgSchema *npschema.Message, opts Options) error {
	ng := numberGenerator{}
	gm := generator.MessageCodeGeneratorMap{
		datatype.Int8:    ng,
		datatype.Int32:   ng,
		datatype.Int64:   ng,
		datatype.Double:  ng,
		datatype.String:  stringGenerator{},
		datatype.Bool:    boolGenerator{},
		datatype.Message: messageGenerator{},
	}
	gm[datatype.Optional] = optionalGenerator{gm}
	gm[datatype.Array] = arrayGenerator{gm}
	gm[datatype.Map] = mapGenerator{gm}

	var wg sync.WaitGroup
	errs := make([]error, 2)

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

	wg.Wait()

	if errs[0] != nil || errs[1] != nil {
		return errors.Join(errs...)
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

func generateMessageHeaderFile(msgSchema *npschema.Message, gm generator.MessageCodeGeneratorMap, opts Options) error {
	info := messageHeaderFileTemplateInfo{
		MessageName:      msgSchema.Name,
		TypeID:           msgSchema.TypeID,
		HasParentMessage: msgSchema.HasParentMessage,
		IncludeGuardName: strcase.ToScreamingSnake(msgSchema.Name) + "_NP_HXX",
		IsInherited:      msgSchema.IsInherited,
	}

	libimp := map[string]struct{}{}
	relimp := map[string]struct{}{}
	for _, field := range msgSchema.DeclaredFields {
		switch field.Type.Kind {
		case datatype.String:
			libimp["string"] = struct{}{}

		case datatype.Map:
			libimp["unordered_map"] = struct{}{}

		case datatype.Optional:
			libimp["optional"] = struct{}{}

		case datatype.Message:
			p, err := resolveSchemaImportPath(field.Type.Schema, msgSchema)
			if err != nil {
				return err
			}
			relimp[p] = struct{}{}

		default:
			continue
		}
	}

	if msgSchema.HasParentMessage {
		info.ParentMessageName = msgSchema.ParentMessage.Name
		p, err := resolveSchemaImportPath(msgSchema.ParentMessage, msgSchema)
		if err != nil {
			return err
		}
		relimp[p] = struct{}{}
	}

	for k := range libimp {
		info.LibraryImports = append(info.LibraryImports, k)
	}
	for k := range relimp {
		info.RelativeImports = append(info.RelativeImports, k)
	}

	for _, f := range msgSchema.InheritedFields {
		g := gm[f.Type.Kind]
		info.ConstructorParameters = append(info.ConstructorParameters, g.ConstructorFieldParameter(f))
	}

	for _, f := range msgSchema.DeclaredFields {
		g := gm[f.Type.Kind]
		info.FieldDeclarationLines = append(info.FieldDeclarationLines, g.FieldDeclaration(f))
		info.ConstructorParameters = append(info.ConstructorParameters, g.ConstructorFieldParameter(f))
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

	op := strings.Replace(msgSchema.SchemaPath, filepath.Base(msgSchema.SchemaPath), fname, 1)
	f, err := os.Create(op)
	if err != nil {
		return err
	}

	defer f.Close()

	err = tmpl.Execute(f, info)
	if err != nil {
		return err
	}

	err = formatCode(op, opts.FormatterPath, opts.FormatterArgs...)
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
		HeaderName:              strings.Replace(fname, extImplFile, extHeaderFile, 1),
		MessageName:             msgSchema.Name,
		HasParentMessage:        msgSchema.HasParentMessage,
		ReadPtrStart:            npHeaderByteSize,
		ConstructorParameters:   nil,
		SuperConstructorArgs:    nil,
		FieldInitializers:       nil,
		InitialWriteBufferSize:  npHeaderByteSize,
		FieldReadCodeFragments:  nil,
		FieldWriteCodeFragments: nil,
	}
	if msgSchema.HasParentMessage {
		info.ParentMessageName = msgSchema.ParentMessage.Name
	}

	for _, f := range msgSchema.InheritedFields {
		g := gm[f.Type.Kind]
		info.ConstructorParameters = append(info.ConstructorParameters, g.ConstructorFieldParameter(f))
		info.SuperConstructorArgs = append(info.SuperConstructorArgs, strcase.ToSnake(f.Name))
	}

	for _, f := range msgSchema.DeclaredFields {
		g := gm[f.Type.Kind]
		info.ConstructorParameters = append(info.ConstructorParameters, g.ConstructorFieldParameter(f))
		info.FieldInitializers = append(info.FieldInitializers, g.ConstructorFieldInitializer(f))
	}

	ctx := generator.CodeContext{}
	for _, f := range msgSchema.AllFields {
		g := gm[f.Type.Kind]
		info.FieldReadCodeFragments = append(info.FieldReadCodeFragments, g.ReadFieldFromBuffer(f, ctx))
	}

	// new context in new block of code
	ctx = generator.CodeContext{}
	for _, f := range msgSchema.AllFields {
		g := gm[f.Type.Kind]
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

	op := strings.Replace(msgSchema.SchemaPath, filepath.Base(msgSchema.SchemaPath), fname, 1)
	f, err := os.Create(op)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, info)
	if err != nil {
		return err
	}

	err = formatCode(op, opts.FormatterPath, opts.FormatterArgs...)
	if err != nil {
		return err
	}

	return nil
}

func generateMessageFactoryHeaderFile(opts Options) error {
	op := filepath.Join(opts.MessageFactoryPath, fileNameMessageFactory+extHeaderFile)
	err := os.WriteFile(op, []byte(messageFactoryHeaderFile), 0644)
	if err != nil {
		return err
	}
	return formatCode(op, opts.FormatterPath, opts.FormatterArgs...)
}

func generateMessageFactoryImplFile(schemas []*npschema.Message, opts Options) error {
	op := filepath.Join(opts.MessageFactoryPath, fileNameMessageFactory+extImplFile)

	info := messageFactoryImplFileTemplateInfo{
		MessageImportPaths: nil,
		MessageSchemas:     nil,
	}

	for _, s := range schemas {
		ip, err := filepath.Rel(filepath.Dir(op), s.SchemaPath)
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

	f, err := os.Create(op)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, info)
	if err != nil {
		return err
	}

	err = formatCode(op, opts.FormatterPath, opts.FormatterArgs...)
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
