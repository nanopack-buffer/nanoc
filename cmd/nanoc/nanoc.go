package main

import (
	"flag"
	"log"
	"nanoc/internal/nanoc"
	"path/filepath"
)

func main() {
	var language string
	var factoryOut string

	flag.StringVar(&language, "language", "", "The language of the generated code.")
	flag.StringVar(&factoryOut, "factory-out", "", "Optionally generate a message factory.")

	flag.Parse()

	if !nanoc.IsLanguageSupported(language) {
		log.Fatalln(language + " is not yet supported by nanoc!")
	}

	fout, err := filepath.Abs(factoryOut)
	if err != nil {
		log.Fatalln("Message factory path is invalid. Received: " + factoryOut)
	}

	inputs := make([]string, flag.NArg())
	for i, p := range flag.Args() {
		abs, err := filepath.Abs(p)
		if err != nil {
			log.Fatalln("Invalid input path encountered: " + p)
		}
		inputs[i] = abs
	}

	opts := nanoc.Options{
		Language:                  nanoc.SupportedLanguage(language),
		MessageFactoryAbsFilePath: fout,
		InputFileAbsolutePaths:    inputs,
	}

	if opts.CodeFormatterPath == "" {
		opts.CodeFormatterPath, opts.CodeFormatterArgs = nanoc.DefaultFormatter(opts.Language)
	}

	err = nanoc.Run(opts)
	if err != nil {
		log.Fatal(err)
	}
}
