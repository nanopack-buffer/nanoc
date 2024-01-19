package swiftgen

import (
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
	"os"
	"path/filepath"
	"strings"
	"text/template"
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

	var t string
	if enumSchema.ValueType.Kind == datatype.String {
		t = "case %v = \"%v\""
	} else {
		t = "case %v = %v"
	}

	for _, m := range enumSchema.Members {
		info.MemberDeclarations = append(info.MemberDeclarations, fmt.Sprintf(t, strcase.ToLowerCamel(m.Name), m.ValueLiteral))
	}

	tmpl, err := template.New(templateNameEnum).Parse(enumTemplate)
	if err != nil {
		return err
	}

	fname := filepath.Base(enumSchema.SchemaPath)
	fname = strcase.ToCamel(strings.TrimSuffix(fname, filepath.Ext(fname))) + extSwift

	op := strings.Replace(enumSchema.SchemaPath, filepath.Base(enumSchema.SchemaPath), fname, 1)
	f, err := os.Create(op)
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
