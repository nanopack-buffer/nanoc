package pathutil

import (
	"nanoc/internal/datatype"
	"os"
	"path/filepath"
	"strings"
)

// ResolveCodeOutputPathForSchema returns the absolute output path to the generated code file for the given schema.
// baseDir is the base directory containing all the schema files currently being compiled (specified by the user),
// outDir is the output directory containing all the generated code files (also specified by the user),
// and outName specifies the file name of the generated code file (e.g. my-message.np.ts).
func ResolveCodeOutputPathForSchema(schema datatype.Schema, baseDir, outDir, outName string) string {
	path := schema.SchemaPathAbsolute()
	fileName := filepath.Base(path)
	path = strings.Replace(path, baseDir, outDir, 1)
	path = strings.Replace(path, fileName, outName, 1)
	return path
}

// CreateOutputFile creates a file at the given absolute path, creating directories if necessary.
func CreateOutputFile(absPath string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(absPath), os.ModePerm)
	if err != nil {
		return nil, err
	}
	return os.Create(absPath)
}
