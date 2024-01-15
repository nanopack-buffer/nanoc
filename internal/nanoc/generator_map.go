package nanoc

import (
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
)

type schemaGeneratorFunc func(schema datatype.Schema, opts Options) error
type messageFactoryGeneratorFunc func(schemas []*npschema.Message, opts Options) error

var schemaGeneratorMap = map[SupportedLanguage]schemaGeneratorFunc{
	LanguageCxx: runCxxSchemaGenerator,
}

var messageFactoryGeneratorMap = map[SupportedLanguage]messageFactoryGeneratorFunc{
	LanguageCxx: runCxxMessageFactoryGenerator,
}
