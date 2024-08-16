package tsgen

import (
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
	"nanoc/internal/symbol"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
)

func outputFileNameForSchema(schema datatype.Schema) string {
	switch s := schema.(type) {
	case *npschema.Enum:
		fname := filepath.Base(s.SchemaPath)
		fname = strcase.ToKebab(strings.TrimSuffix(fname, symbol.SchemaFileExt)) + extTsFile
		return fname

	case *npschema.Message:
		fname := filepath.Base(s.SchemaPath)
		fname = strcase.ToKebab(strings.TrimSuffix(fname, symbol.SchemaFileExt)) + extTsFile
		return fname

	case *npschema.Service:
		return strcase.ToKebab(s.Name+"Service") + extTsFile

	default:
		return ""
	}
}
