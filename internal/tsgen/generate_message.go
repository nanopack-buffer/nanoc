package tsgen

import (
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"os"
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
	if !msgSchema.IsInherited {
		return generateMessageClass(msgSchema, opts)
	}

	var wg sync.WaitGroup
	errs := make([]error, 2)

	wg.Add(1)
	go func() {
		errs[0] = generateMessageClass(msgSchema, opts)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errs[1] = generateMessageClassFactory(msgSchema, opts)
		wg.Done()
	}()

	wg.Wait()

	if errs[0] != nil || errs[1] != nil {
		return errors.Join(errs...)
	}

	return nil
}

func GenerateMessageFactory(schemas []*npschema.Message, opts Options) error {
	info := messageFactoryTemplateInfo{
		Schemas: schemas,
	}

	op := filepath.Join(opts.MessageFactoryPath, fileNameMessageFactoryFile+extTsFile)

	for _, s := range schemas {
		p, err := resolveImportPath(s.SchemaPath, op)
		if err != nil {
			return err
		}
		info.MessageImports = append(info.MessageImports, fmt.Sprintf("import { %v } from \"%v\";", s.Name, p))
	}

	tmpl, err := template.New(templateNameMessageFactory).Parse(messageFactoryTemplate)
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

func generateMessageClass(msgSchema *npschema.Message, opts Options) error {
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
	gm[datatype.Enum] = enumGenerator{gm}

	// the message header byte size includes 4 bytes for the type ID and 4 bytes to store the byte size of each field
	npHeaderByteSize := (len(msgSchema.AllFields) + 1) * 4

	info := messageClassTemplateInfo{
		Schema:                 msgSchema,
		ReadPtrStart:           npHeaderByteSize,
		InitialWriteBufferSize: npHeaderByteSize,
	}

	for _, s := range msgSchema.ImportedTypes {
		p, err := resolveSchemaImportPath(s, msgSchema)
		if err != nil {
			return err
		}

		switch s := s.(type) {
		case *npschema.Message:
			info.ExternalImports = append(info.ExternalImports, fmt.Sprintf("import { %v } from \"%v\";", s.Name, p))
		case *npschema.Enum:
			info.ExternalImports = append(info.ExternalImports, fmt.Sprintf("import type { T%v } from \"%v\";", s.Name, p))
		}
	}

	for _, f := range msgSchema.InheritedFields {
		g := gm[f.Type.Kind]
		c := strcase.ToLowerCamel(f.Name)
		info.ConstructorParameters = append(info.ConstructorParameters, g.ConstructorFieldParameter(f))
		info.ConstructorArgs = append(info.ConstructorArgs, c)
		info.SuperConstructorArgs = append(info.SuperConstructorArgs, c)
	}
	for _, f := range msgSchema.DeclaredFields {
		g := gm[f.Type.Kind]
		c := strcase.ToLowerCamel(f.Name)
		info.ConstructorArgs = append(info.ConstructorArgs, c)
		info.ConstructorParameters = append(info.ConstructorParameters, g.ConstructorFieldParameter(f))
	}

	ctx := generator.NewCodeContext()
	for _, f := range msgSchema.AllFields {
		g := gm[f.Type.Kind]
		info.FieldReadCodeFragments = append(info.FieldReadCodeFragments, g.ReadFieldFromBuffer(f, ctx))
	}

	// new context in new block of code
	ctx = generator.NewCodeContext()
	for _, f := range msgSchema.AllFields {
		g := gm[f.Type.Kind]
		info.FieldWriteCodeFragments = append(info.FieldWriteCodeFragments, g.WriteFieldToBuffer(f, ctx))
	}

	tmpl, err := template.New(templateNameMessageClass).
		Funcs(template.FuncMap{
			"join": strings.Join,
		}).
		Parse(messageClassTemplate)

	if err != nil {
		return err
	}

	fname := filepath.Base(msgSchema.SchemaPath)
	fname = strcase.ToKebab(strings.TrimSuffix(fname, filepath.Ext(fname))) + extTsFile

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

func generateMessageClassFactory(msgSchema *npschema.Message, opts Options) error {
	info := messageClassFactoryTemplateInfo{
		Schema:             msgSchema,
		MessageClassImport: "./" + strcase.ToKebab(msgSchema.Name) + extImport,
		MessageImports:     nil,
	}

	for _, f := range msgSchema.ChildMessages {
		p, err := resolveSchemaImportPath(f, msgSchema)
		if err != nil {
			return err
		}
		info.MessageImports = append(info.MessageImports, fmt.Sprintf("import { %v } from \"%v\";", f.Name, p))
	}

	tmpl, err := template.New(templateNameMessageClassFactory).Parse(messageClassFactoryTemplate)
	if err != nil {
		return err
	}

	kb := strcase.ToKebab(msgSchema.Name)

	op := strings.Replace(msgSchema.SchemaPath, filepath.Base(msgSchema.SchemaPath), fmt.Sprintf("make-%v%v", kb, extTsFile), 1)
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
