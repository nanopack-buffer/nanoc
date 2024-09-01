package cxxgen

import (
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
	info := enumHeaderFileInfo{
		Namespace:        strings.Join(opts.Namespaces, cxxSymbolMemberOf),
		Schema:           enumSchema,
		BackingTypeName:  "",
		MemberNames:      nil,
		IncludeGuardName: fmt.Sprintf("%v_ENUM_NP_HXX", strcase.ToScreamingSnake(enumSchema.Name)),
	}

	switch enumSchema.ValueType.Kind {
	case datatype.Int8:
		info.BackingTypeName = "int8_t"
	case datatype.Int32:
		info.BackingTypeName = "int32_t"
	case datatype.Int64:
		info.BackingTypeName = "int64_t"
	case datatype.Double:
		info.BackingTypeName = "double"
	case datatype.String:
		info.BackingTypeName = "std::string_view"
	default:
		info.BackingTypeName = ""
	}

	for _, m := range enumSchema.Members {
		info.MemberNames = append(info.MemberNames, strcase.ToScreamingSnake(m.Name))
	}

	tmpl, err := template.New(templateNameEnumHeaderFile).Parse(enumHeaderFile)
	if err != nil {
		return err
	}

	fname := filepath.Base(enumSchema.SchemaPath)
	fname = strcase.ToSnake(strings.TrimSuffix(fname, symbol.SchemaFileExt)) + extHeaderFile
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
