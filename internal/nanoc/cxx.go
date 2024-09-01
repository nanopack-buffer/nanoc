package nanoc

import (
	"errors"
	"nanoc/internal/cxxgen"
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
	"reflect"
	"strings"
)

func runCxxSchemaGenerator(schema datatype.Schema, opts Options) error {
	o := cxxgen.Options{
		FormatterPath:       opts.CodeFormatterPath,
		FormatterArgs:       opts.CodeFormatterArgs,
		MessageFactoryPath:  opts.MessageFactoryAbsFilePath,
		Namespaces:          strings.Split(opts.Namespace, "."),
		BaseDirectoryPath:   opts.BaseDirectoryAbs,
		OutputDirectoryPath: opts.OutputDirectoryAbs,
	}

	switch s := schema.(type) {
	case *npschema.Message:
		return cxxgen.GenerateMessageClass(s, o)

	case *npschema.Enum:
		return cxxgen.GenerateEnum(s, o)

	case *npschema.Service:
		return cxxgen.GenerateService(s, o)

	default:
		return errors.New("unexpected error. Unsupported schema type " + reflect.TypeOf(schema).Name())
	}
}

func runCxxMessageFactoryGenerator(schema []*npschema.Message, opts Options) error {
	return cxxgen.GenerateMessageFactory(schema, cxxgen.Options{
		Namespaces:          strings.Split(opts.Namespace, "."),
		FormatterPath:       opts.CodeFormatterPath,
		FormatterArgs:       opts.CodeFormatterArgs,
		MessageFactoryPath:  opts.MessageFactoryAbsFilePath,
		BaseDirectoryPath:   opts.BaseDirectoryAbs,
		OutputDirectoryPath: opts.OutputDirectoryAbs,
	})
}
