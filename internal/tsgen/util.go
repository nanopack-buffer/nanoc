package tsgen

import (
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"os/exec"
	"path/filepath"
	"strings"
)

func resolveImportPath(toPath string, fromPath string) (string, error) {
	p, err := filepath.Rel(filepath.Dir(toPath), fromPath)
	if err != nil {
		return "", err
	}

	fname := filepath.Base(toPath)
	fname = strcase.ToKebab(strings.TrimSuffix(fname, filepath.Ext(fname))) + extImport

	p = strings.Replace(p, filepath.Base(p), fname, 1)
	if !strings.HasPrefix(p, ".") || !strings.HasPrefix(p, "/") {
		p = "./" + p
	}

	return p, nil
}

func resolveSchemaImportPath(toSchema datatype.Schema, fromSchema datatype.Schema) (string, error) {
	return resolveImportPath(toSchema.SchemaPathAbsolute(), fromSchema.SchemaPathAbsolute())
}

func formatCode(path string, formatter string, args ...string) error {
	args = append(args, path)
	cmd := exec.Command(formatter, args...)
	return cmd.Run()
}
