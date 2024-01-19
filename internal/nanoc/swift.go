package nanoc

import (
	"errors"
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
	"nanoc/internal/swiftgen"
	"reflect"
)

func runSwiftSchemaGenerator(schema datatype.Schema, opts Options) error {
	o := swiftgen.Options{
		FormatterPath:      opts.CodeFormatterPath,
		FormatterArgs:      opts.CodeFormatterArgs,
		MessageFactoryPath: opts.MessageFactoryAbsFilePath,
	}

	switch s := schema.(type) {
	case *npschema.Message:
		return swiftgen.GenerateMessageClass(s, o)

	case *npschema.Enum:
		return swiftgen.GenerateEnum(s, o)

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
