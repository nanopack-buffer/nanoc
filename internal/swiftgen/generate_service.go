package swiftgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"nanoc/internal/pathutil"
	"os"
	"strings"
	"text/template"
	"unicode/utf8"
)

func GenerateService(serviceSchema *npschema.Service, opts Options) error {
	info := serviceTemplateInfo{
		Schema: serviceSchema,
	}

	ng := numberGenerator{}
	gm := generator.MessageCodeGeneratorMap{
		datatype.Int8:    ng,
		datatype.Int32:   ng,
		datatype.Int64:   ng,
		datatype.UInt8:   ng,
		datatype.UInt32:  ng,
		datatype.UInt64:  ng,
		datatype.Double:  ng,
		datatype.String:  stringGenerator{},
		datatype.Bool:    boolGenerator{},
		datatype.Message: messageGenerator{},
		datatype.Any:     anyGenerator{},
	}
	gm[datatype.Optional] = optionalGenerator{gm}
	gm[datatype.Array] = arrayGenerator{gm}
	gm[datatype.Map] = mapGenerator{gm}
	gm[datatype.Enum] = enumGenerator{gm}

	// a dummy code context that will be passed to generator
	// this will not be used
	dummyCtx := generator.CodeContext{}

	funcs := template.FuncMap{
		"stringByteSize": func(s string) int {
			return utf8.RuneCountInString(s)
		},
		"generateReadParamCode": func(fn *npschema.DeclaredFunction) (string, error) {
			var sb strings.Builder
			for _, param := range fn.Parameters {
				g := gm[param.Type.Kind]
				sb.WriteString(g.ReadValueFromBuffer(param.Type, param.Name, dummyCtx))
				sb.WriteRune('\n')
			}
			return sb.String(), nil
		},
		"generateWriteParamCode": func(fn *npschema.DeclaredFunction) (string, error) {
			var sb strings.Builder
			for _, param := range fn.Parameters {
				g := gm[param.Type.Kind]
				sb.WriteString(g.WriteVariableToBuffer(param.Type, param.Name, dummyCtx))
				sb.WriteRune('\n')
			}
			return sb.String(), nil
		},
		"generateReadResultCode": func(fn *npschema.DeclaredFunction) string {
			g := gm[fn.ReturnType.Kind]
			return g.ReadValueFromBuffer(*fn.ReturnType, "result", dummyCtx)
		},
		"generateWriteResultCode": func(fn *npschema.DeclaredFunction) string {
			g := gm[fn.ReturnType.Kind]
			return g.WriteVariableToBuffer(*fn.ReturnType, "result", dummyCtx)
		},
		"typeDeclaration": func(t datatype.DataType) string {
			g := gm[t.Kind]
			return g.TypeDeclaration(t)
		},
	}

	tmpl, err := template.New(templateNameServiceServer).Funcs(funcs).Parse(serviceTemplate)
	if err != nil {
		return err
	}

	fname := fmt.Sprintf("%vServiceServer.swift", serviceSchema.Name)
	outPath := pathutil.ResolveCodeOutputPathForSchema(serviceSchema, opts.BaseDirectoryPath, opts.OutputDirectoryPath, fname)
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, info)
	if err != nil {
		return err
	}

	err = formatCode(outPath, opts.FormatterPath, opts.FormatterArgs...)
	if err != nil {
		return err
	}

	return nil
}
