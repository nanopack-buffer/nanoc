package main

import (
	"flag"
	"log"
	"nanoc/internal/nanoc"
	"os"
	"path/filepath"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var language string
	var factoryOut string
	var namespace string
	var baseDir string
	var outputDir string

	flag.StringVar(&language, "language", "", "The language of the generated code.")
	flag.StringVar(&factoryOut, "factory-out", "", "Optionally generate a message factory.")
	flag.StringVar(&namespace, "namespace", "", "Optionally put the generated code under this namespace. Use dot notation for nested namespaces, for example My.Namespace")
	flag.StringVar(&baseDir, "basedir", "", "The base directory in which schema files are stored. All schema files that are compiled together must be placed under this directory (subdirectories are fine and will be preseved in the output directory). Default is the current directory.")
	flag.StringVar(&outputDir, "outdir", "", "The directory in which generated code will be placed. Subdirectories of the base directory are preseved. Default is the current working directory.")

	flag.Parse()

	if !nanoc.IsLanguageSupported(language) {
		log.Fatalln(language + " is not yet supported by nanoc!")
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
		MessageFactoryAbsFilePath: "",
		InputFileAbsolutePaths:    inputs,
		Namespace:                 namespace,
	}

	if baseDir != "" {
		baseDirAbs, err := filepath.Abs(baseDir)
		if err != nil {
			log.Fatalln("Unable to resolve the path of the specified base directory. Make sure the directory exists!")
		}
		opts.BaseDirectoryAbs = baseDirAbs
	} else {
		opts.BaseDirectoryAbs = cwd
	}

	if outputDir != "" {
		outputDirAbs, err := filepath.Abs(outputDir)
		if err != nil {
			log.Fatalln("Unable to resolve the path of the specified output directory. Make sure the directory exists!")
		}
		opts.OutputDirectoryAbs = outputDirAbs
	} else {
		opts.OutputDirectoryAbs = cwd
	}

	if factoryOut != "" {
		p, err := filepath.Abs(filepath.Join(opts.OutputDirectoryAbs, factoryOut))
		if err != nil {
			log.Fatalln("Message factory path is invalid. Received: " + factoryOut)
		}
		opts.MessageFactoryAbsFilePath = p
	}

	if opts.CodeFormatterPath == "" {
		opts.CodeFormatterPath, opts.CodeFormatterArgs = nanoc.DefaultFormatter(opts.Language)
	}

	err = nanoc.Run(opts)
	if err != nil {
		log.Fatal(err)
	}
}
