package tsgen

import (
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
		Schema: enumSchema,
	}

	for _, m := range enumSchema.Members {
		if enumSchema.ValueType.Kind == datatype.String {
			info.MemberDeclarations = append(info.MemberDeclarations, fmt.Sprintf("%v: \"%v\"", m.Name, m.ValueLiteral))
		} else {
			info.MemberDeclarations = append(info.MemberDeclarations, fmt.Sprintf("%v: %v", m.Name, m.ValueLiteral))
		}
	}

	tmpl, err := template.New(templateNameEnum).Parse(enumTemplate)
	if err != nil {
		return err
	}

	fname := filepath.Base(enumSchema.SchemaPath)
	fname = strcase.ToKebab(strings.TrimSuffix(fname, filepath.Ext(fname))) + extTsFile

	op := strings.Replace(enumSchema.SchemaPath, filepath.Base(enumSchema.SchemaPath), fname, 1)
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
