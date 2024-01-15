package nanoc

import (
	"errors"
	"nanoc/internal/cxxgenerator"
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
	"nanoc/internal/parser"
	"nanoc/internal/resolver"
	"reflect"
	"sync"
)

type generatorFunc func(schema datatype.Schema, opts Options) error

func Run(opts Options) error {
	var wg sync.WaitGroup
	sc := len(opts.InputFileAbsolutePaths)

	errs := make([]error, 0, sc)

	partialSchemas := make([]datatype.PartialSchema, 0, sc)
	mu := &sync.Mutex{}
	for _, p := range opts.InputFileAbsolutePaths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			sp, err := parser.ParseSchema(path)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			} else {
				mu.Lock()
				partialSchemas = append(partialSchemas, sp)
				mu.Unlock()
			}
		}(p)
	}
	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	schemas, err := resolver.Resolve(partialSchemas)
	if err != nil {
		return err
	}

	var generator generatorFunc
	switch opts.Language {
	case LanguageCxx:
		generator = runCxxGenerator
	}

	for _, s := range schemas {
		s := s
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := generator(s, opts)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func runCxxGenerator(schema datatype.Schema, opts Options) error {
	switch s := schema.(type) {
	case *npschema.Message:
		return cxxgenerator.GenerateMessageClass(s, cxxgenerator.Options{
			FormatterPath: opts.CodeFormatterPath,
			FormatterArgs: opts.CodeFormatterArgs,
		})

	default:
		return errors.New("unexpected error. Unsupported schema type " + reflect.TypeOf(schema).Name())
	}
}
