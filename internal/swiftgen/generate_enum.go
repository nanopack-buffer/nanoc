package swiftgen

import (
	"errors"
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
	"nanoc/internal/pathutil"
	"nanoc/internal/symbol"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

func GenerateEnum(enumSchema *npschema.Enum, opts Options) error {
	info := enumTemplateInfo{
		EnumName: enumSchema.Name,
	}

	switch enumSchema.ValueType.Kind {
	case datatype.Int8:
		info.SwiftValueTypeName = "Int8"
	case datatype.Int32:
		info.SwiftValueTypeName = "Int32"
	case datatype.Int64:
		info.SwiftValueTypeName = "Int64"
	case datatype.String:
		info.SwiftValueTypeName = "String"
	default:
		return errors.New("not gonna happen")
	}

	for _, m := range enumSchema.Members {
		info.MemberDeclarations = append(info.MemberDeclarations, fmt.Sprintf("case %v = %v", strcase.ToLowerCamel(m.Name), m.ValueLiteral))
	}

	tmpl, err := template.New(templateNameEnum).Parse(enumTemplate)
	if err != nil {
		return err
	}

	fname := filepath.Base(enumSchema.SchemaPath)
	fname = strcase.ToCamel(strings.TrimSuffix(fname, symbol.SchemaFileExt)) + extSwift

	op := pathutil.ResolveCodeOutputPathForSchema(enumSchema, opts.BaseDirectoryPath, opts.OutputDirectoryPath, fname)
	f, err := pathutil.CreateOutputFile(op)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, info)
	if err != nil {
		return nil
	}

	err = formatCode(op, opts.FormatterPath, opts.FormatterArgs...)
	if err != nil {
		return err
	}

	return nil
}
