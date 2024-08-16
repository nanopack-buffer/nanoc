package tsgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
	"nanoc/internal/pathutil"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
)

func resolveImportPaths(importedTypes []datatype.DataType, fromSchema datatype.Schema, basedir, outdir, factoryPath string) ([]string, error) {
	msgFactoryImported := false

	importPaths := []string{}

	for _, t := range importedTypes {
		if t.Schema != nil {
			p, err := resolveSchemaImportPath(t.Schema, fromSchema, basedir, outdir)
			if err != nil {
				return nil, err
			}

			switch t.Kind {
			case datatype.Message:
				s := t.Schema.(*npschema.Message)
				if s.IsInherited {
					factoryPath := strings.Replace(p, filepath.Base(p), fmt.Sprintf("make-%v%v", strcase.ToKebab(s.Name), extImport), 1)
					importPaths = append(
						importPaths,
						fmt.Sprintf("import { %v } from \"%v\";", s.Name, p),
						fmt.Sprintf("import { make%v } from \"%v\";", s.Name, factoryPath),
					)
				} else {
					importPaths = append(importPaths, fmt.Sprintf("import { %v } from \"%v\";", s.Name, p))
				}

			case datatype.Enum:
				importPaths = append(importPaths, fmt.Sprintf("import type { T%v } from \"%v\";", t.Schema.(*npschema.Enum).Name, p))

			default:
				break
			}
		} else if t.Kind == datatype.Message && !msgFactoryImported {
			p, err := resolveMessageFactoryImportPath(factoryPath, fromSchema, basedir, outdir)
			if err != nil {
				return nil, err
			}
			importPaths = append(importPaths, "import type { NanoPackMessage } from \"nanopack\";", fmt.Sprintf("import { makeNanoPackMessage } from \"%v\";", p))
			msgFactoryImported = true
		}
	}

	return importPaths, nil
}

func resolveImportPath(toPath string, fromPath string) (string, error) {
	p, err := filepath.Rel(filepath.Dir(fromPath), toPath)
	if err != nil {
		return "", err
	}

	fname := filepath.Base(toPath)
	if strings.HasSuffix(fname, extTsFile) {
		fname = strings.Replace(fname, extTsFile, extImport, 1)
	}

	p = strings.Replace(p, filepath.Base(p), fname, 1)
	if !strings.HasPrefix(p, ".") && !strings.HasPrefix(p, "/") {
		p = "./" + p
	}

	return p, nil
}

func resolveSchemaImportPath(toSchema datatype.Schema, fromSchema datatype.Schema, basedir string, outdir string) (string, error) {
	dest := pathutil.ResolveCodeOutputPathForSchema(toSchema, basedir, outdir, outputFileNameForSchema(toSchema))
	src := pathutil.ResolveCodeOutputPathForSchema(fromSchema, basedir, outdir, outputFileNameForSchema(fromSchema))
	return resolveImportPath(dest, src)
}

func resolveMessageFactoryImportPath(factoryPath string, fromSchema datatype.Schema, basedir string, outdir string) (string, error) {
	p := filepath.Join(factoryPath, fileNameMessageFactoryFile)
	src := pathutil.ResolveCodeOutputPathForSchema(fromSchema, basedir, outdir, outputFileNameForSchema(fromSchema))
	return resolveImportPath(p, src)
}
