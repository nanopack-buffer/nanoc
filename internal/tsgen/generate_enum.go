package tsgen

import (
	"fmt"
	"nanoc/internal/npschema"
	"nanoc/internal/pathutil"
	"text/template"
)

func GenerateEnum(enumSchema *npschema.Enum, opts Options) error {
	info := enumTemplateInfo{
		Schema: enumSchema,
	}

	for _, m := range enumSchema.Members {
		info.MemberDeclarations = append(info.MemberDeclarations, fmt.Sprintf("%v: %v", m.Name, m.ValueLiteral))
	}

	tmpl, err := template.New(templateNameEnum).Parse(enumTemplate)
	if err != nil {
		return err
	}

	fname := outputFileNameForSchema(enumSchema)
	op := pathutil.ResolveCodeOutputPathForSchema(enumSchema, opts.BaseDirectoryPath, opts.OutputDirectoryPath, fname)
	f, err := pathutil.CreateOutputFile(op)
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
