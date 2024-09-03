package cxxgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"nanoc/internal/pathutil"
	"path/filepath"
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

	libimp := map[string]struct{}{}
	relimp := map[string]struct{}{}
	for _, t := range serviceSchema.ImportedTypes {
		datatype.TraverseTypeTree(&t, func(t *datatype.DataType) error {
			switch t.Kind {
			case datatype.String:
				libimp["string"] = struct{}{}

			case datatype.Map:
				libimp["unordered_map"] = struct{}{}

			case datatype.Optional:
				libimp["optional"] = struct{}{}

			case datatype.Any:
				libimp["nanopack/any.hxx"] = struct{}{}

			case datatype.Message:
				if t.Schema == nil {
					libimp["memory"] = struct{}{}
					libimp["nanopack/message.hxx"] = struct{}{}
				} else {
					if t.Schema.(*npschema.Message).IsInherited {
						libimp["memory"] = struct{}{}
					}

					p, err := resolveSchemaImportPath(t.Schema, serviceSchema, opts.BaseDirectoryPath, opts.OutputDirectoryPath)
					if err != nil {
						return err
					}
					relimp[p] = struct{}{}

					// if this type is inherited (polymorphic) and is imported because it is used by one of the fields
					// then its factory needs to be imported as well,
					// because it will be used when reading fields that use this polymorphic type to instantiate the correct type.
					if ms, ok := t.Schema.(*npschema.Message); ok && ms.IsInherited {
						header := fmt.Sprintf("make_%v%v", strcase.ToSnake(ms.Name), extHeaderFile)
						relimp[strings.Replace(p, filepath.Base(p), header, 1)] = struct{}{}
					}
				}

			case datatype.Array:
				libimp["vector"] = struct{}{}

			default:
				break
			}

			return nil
		})
	}

	for p, _ := range libimp {
		info.LibraryImports = append(info.LibraryImports, p)
	}
	for p, _ := range relimp {
		info.RelativeImports = append(info.RelativeImports, p)
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
