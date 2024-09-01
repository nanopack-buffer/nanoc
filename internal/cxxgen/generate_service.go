package cxxgen

import (
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"nanoc/internal/pathutil"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/iancoleman/strcase"
)

func GenerateService(serviceSchema *npschema.Service, opts Options) error {
	err := generateServiceHeader(serviceSchema, opts)
	if err != nil {
		return err
	}

	err = generateServiceImpl(serviceSchema, opts)
	if err != nil {
		return err
	}

	return nil
}

func generateServiceHeader(serviceSchema *npschema.Service, opts Options) error {
	info := serviceHeaderFileInfo{
		Schema:           serviceSchema,
		Namespace:        strings.Join(opts.Namespaces, cxxSymbolMemberOf),
		IncludeGuardName: strcase.ToScreamingSnake(serviceSchema.Name) + "_SERVICE_NP_HXX",
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

	ctx := generator.NewCodeContext()

	funcs := template.FuncMap{
		"snake": func(s string) string {
			return strcase.ToSnake(s)
		},
		"stringByteSize": func(s string) int {
			return utf8.RuneCountInString(s)
		},
		"typeDeclaration": func(t datatype.DataType) string {
			g := gm[t.Kind]
			return g.TypeDeclaration(t)
		},
		"generateReadParamCode": func(fn *npschema.DeclaredFunction) string {
			var sb strings.Builder
			for _, param := range fn.Parameters {
				g := gm[param.Type.Kind]
				sb.WriteString(g.ReadValueFromBuffer(param.Type, strcase.ToSnake(param.Name), ctx))
				sb.WriteRune('\n')
			}
			return sb.String()
		},
		"generateWriteParamCode": func(fn *npschema.DeclaredFunction) string {
			var sb strings.Builder
			for _, param := range fn.Parameters {
				g := gm[param.Type.Kind]
				sb.WriteString(g.WriteVariableToBuffer(param.Type, strcase.ToSnake(param.Name), ctx))
				sb.WriteRune('\n')
			}
			return sb.String()
		},
		"generateReadResultCode": func(fn *npschema.DeclaredFunction) string {
			g := gm[fn.ReturnType.Kind]
			return g.ReadValueFromBuffer(*fn.ReturnType, "result", ctx)
		},
		"generateWriteResultCode": func(fn *npschema.DeclaredFunction) string {
			g := gm[fn.ReturnType.Kind]
			return g.WriteVariableToBuffer(*fn.ReturnType, "result", ctx)
		},
	}

	tmpl, err := template.New(templateNameServiceHeaderFile).Funcs(funcs).Parse(serviceHeaderFile)
	if err != nil {
		return err
	}

	fname := strcase.ToSnake(serviceSchema.Name) + "_service" + extHeaderFile
	outPath := pathutil.ResolveCodeOutputPathForSchema(serviceSchema, opts.BaseDirectoryPath, opts.OutputDirectoryPath, fname)
	f, err := pathutil.CreateOutputFile(outPath)
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

func generateServiceImpl(serviceSchema *npschema.Service, opts Options) error {
	info := serviceImplFileInfo{
		Schema:     serviceSchema,
		HeaderName: strcase.ToSnake(serviceSchema.Name) + "_service" + extHeaderFile,
		Namespace:  strings.Join(opts.Namespaces, cxxSymbolMemberOf),
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

	ctx := generator.NewCodeContext()

	funcs := template.FuncMap{
		"snake": func(s string) string {
			return strcase.ToSnake(s)
		},
		"stringByteSize": func(s string) int {
			return utf8.RuneCountInString(s)
		},
		"typeDeclaration": func(t datatype.DataType) string {
			g := gm[t.Kind]
			return g.TypeDeclaration(t)
		},
		"generateReadParamCode": func(fn *npschema.DeclaredFunction) string {
			var sb strings.Builder
			for _, param := range fn.Parameters {
				g := gm[param.Type.Kind]
				sb.WriteString(g.ReadValueFromBuffer(param.Type, strcase.ToSnake(param.Name), ctx))
				sb.WriteRune('\n')
			}
			return sb.String()
		},
		"generateWriteParamCode": func(fn *npschema.DeclaredFunction) string {
			var sb strings.Builder
			for _, param := range fn.Parameters {
				g := gm[param.Type.Kind]
				sb.WriteString(g.WriteVariableToBuffer(param.Type, strcase.ToSnake(param.Name), ctx))
				sb.WriteRune('\n')
			}
			return sb.String()
		},
		"generateReadResultCode": func(fn *npschema.DeclaredFunction) string {
			g := gm[fn.ReturnType.Kind]
			return g.ReadValueFromBuffer(*fn.ReturnType, "result", ctx)
		},
		"generateWriteResultCode": func(fn *npschema.DeclaredFunction) string {
			g := gm[fn.ReturnType.Kind]
			return g.WriteVariableToBuffer(*fn.ReturnType, "result", ctx)
		},
		"isTriviallyCopyable": func(t datatype.DataType) bool {
			return isTriviallyCopiable(t)
		},
	}

	tmpl, err := template.New(templateNameServiceImplFile).Funcs(funcs).Parse(serviceImplFile)
	if err != nil {
		return err
	}

	fname := strcase.ToSnake(serviceSchema.Name) + "_service" + extImplFile
	outPath := pathutil.ResolveCodeOutputPathForSchema(serviceSchema, opts.BaseDirectoryPath, opts.OutputDirectoryPath, fname)
	f, err := pathutil.CreateOutputFile(outPath)
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