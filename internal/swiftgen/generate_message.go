package swiftgen

import "nanoc/internal/npschema"

// Options are parameters that can be tweaked to alter codegen.
type Options struct {
	FormatterPath string
	FormatterArgs []string

	// The absolute path to the directory where the factory file should be put in
	// This is an empty string when it is not requested.
	MessageFactoryPath string
}

func GenerateMessageClass(msgSchema *npschema.Message, opts Options) error {
	return nil
}
