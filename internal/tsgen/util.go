package tsgen

import (
	"nanoc/internal/datatype"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
)

func resolveImportPath(toPath string, fromPath string) (string, error) {
	p, err := filepath.Rel(filepath.Dir(fromPath), toPath)
	if err != nil {
		return "", err
	}

	fname := filepath.Base(toPath)
	fname = strcase.ToKebab(strings.TrimSuffix(fname, filepath.Ext(fname))) + extImport

	p = strings.Replace(p, filepath.Base(p), fname, 1)
	if !strings.HasPrefix(p, ".") && !strings.HasPrefix(p, "/") {
		p = "./" + p
	}

	return p, nil
}

func resolveSchemaImportPath(toSchema datatype.Schema, fromSchema datatype.Schema) (string, error) {
	return resolveImportPath(toSchema.SchemaPathAbsolute(), fromSchema.SchemaPathAbsolute())
}

func resolveMessageFactoryImportPath(factoryPath string, fromSchema datatype.Schema) (string, error) {
	p := filepath.Join(factoryPath, fileNameMessageFactoryFile)
	return resolveImportPath(p, fromSchema.SchemaPathAbsolute())
}

func formatCode(path string, formatter string, args ...string) error {
	args = append(args, path)
	cmd := exec.Command(formatter, args...)
	return cmd.Run()
}
