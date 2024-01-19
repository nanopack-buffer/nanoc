package swiftgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	// the message header byte size includes 4 bytes for the type ID and 4 bytes to store the byte size of each field
	npHeaderByteSize := (len(msgSchema.AllFields) + 1) * 4

	info := messageClassTemplateInfo{
		Schema:                 msgSchema,
		ReadPtrStart:           npHeaderByteSize,
		InitialWriteBufferSize: npHeaderByteSize,
	}

	for _, f := range msgSchema.InheritedFields {
		g := gm[f.Type.Kind]
		c := strcase.ToLowerCamel(f.Name)
		info.FieldInitializers = append(info.FieldInitializers, g.FieldInitializer(f))
		info.ConstructorParameters = append(info.ConstructorParameters, g.ConstructorFieldParameter(f))
		info.SuperConstructorArgs = append(info.SuperConstructorArgs, fmt.Sprintf("%v: %v", c, c))
	}
	for _, f := range msgSchema.DeclaredFields {
		g := gm[f.Type.Kind]
		info.FieldDeclarationLines = append(info.FieldDeclarationLines, g.FieldDeclaration(f))
		info.ConstructorParameters = append(info.ConstructorParameters, g.ConstructorFieldParameter(f))
		info.FieldInitializers = append(info.FieldInitializers, g.FieldInitializer(f))
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
	fname = strcase.ToCamel(strings.TrimSuffix(fname, filepath.Ext(fname))) + extSwift

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

func GenerateMessageFactory(schemas []*npschema.Message, opts Options) error {
	info := messageFactoryTemplateInfo{
		Schemas: schemas,
	}

	tmpl, err := template.New(templateNameMessageFactory).
		Funcs(template.FuncMap{
			"join": strings.Join,
		}).
		Parse(messageFactoryTemplate)

	if err != nil {
		return err
	}

	op := filepath.Join(opts.MessageFactoryPath, fileNameMessageFactoryFile+extSwift)
	f, err := os.Create(op)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, info)
	if err != nil {
		return err
	}

	return nil
}

func formatCode(path string, formatter string, args ...string) error {
	args = append(args, path)
	cmd := exec.Command(formatter, args...)
	return cmd.Run()
}
