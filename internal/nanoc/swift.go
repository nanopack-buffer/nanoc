package nanoc

import (
	"errors"
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
	"nanoc/internal/swiftgen"
	"reflect"
)

func runSwiftSchemaGenerator(schema datatype.Schema, opts Options) error {
	switch s := schema.(type) {
	case *npschema.Message:
		return swiftgen.GenerateMessageClass(s, swiftgen.Options{
			FormatterPath:      opts.CodeFormatterPath,
			FormatterArgs:      opts.CodeFormatterArgs,
			MessageFactoryPath: opts.MessageFactoryAbsFilePath,
		})
	default:
		return errors.New("unexpected error. Unsupported schema type " + reflect.TypeOf(schema).Name())
	}
}

func runSwiftMessageFactoryGenerator(schemas []*npschema.Message, opts Options) error {
	return swiftgen.GenerateMessageFactory(schemas, swiftgen.Options{
		FormatterPath:      opts.CodeFormatterPath,
		FormatterArgs:      opts.CodeFormatterArgs,
		MessageFactoryPath: opts.MessageFactoryAbsFilePath,
	})
}
