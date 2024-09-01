package cxxgen

import (
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
	"nanoc/internal/symbol"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
)

// outputHeaderFileNameForSchema returns the name of the generated header file for the given schema
func outputHeaderFileNameForSchema(schema datatype.Schema) string {
	switch s := schema.(type) {
	case *npschema.Enum:
		fname := filepath.Base(s.SchemaPath)
		fname = strcase.ToSnake(strings.TrimSuffix(fname, symbol.SchemaFileExt)) + extHeaderFile
		return fname

	case *npschema.Message:
		fname := filepath.Base(s.SchemaPath)
		fname = strcase.ToSnake(strings.TrimSuffix(fname, symbol.SchemaFileExt)) + extHeaderFile
		return fname

	case *npschema.Service:
		return strcase.ToSnake(s.Name) + "_service" + extHeaderFile

	default:
		return ""
	}
}
