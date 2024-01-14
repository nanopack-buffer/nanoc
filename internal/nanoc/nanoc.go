package nanoc

import (
	"errors"
	"nanoc/internal/cxxgenerator"
	"nanoc/internal/npschema"
	"nanoc/internal/parser"
	"nanoc/internal/resolver"
	"reflect"
	"sync"
)

func Run(opts Options) error {
	var wg sync.WaitGroup
	sc := len(opts.InputFileAbsolutePaths)

	errs := make([]error, sc)

	partialSchemas := make([]npschema.PartialSchema, sc)
	mu := &sync.Mutex{}
	for _, p := range opts.InputFileAbsolutePaths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			sp, err := parser.Parse(path)
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

	var generator func(schema npschema.Schema) error
	switch opts.Language {
	case LanguageCxx:
		generator = runCxxGenerator
	}

	for _, s := range schemas {
		s := s
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := generator(s)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}()
	}

	return nil
}

func runCxxGenerator(schema npschema.Schema) error {
	switch s := schema.(type) {
	case npschema.Message:
		return cxxgenerator.GenerateMessageClass(s)

	default:
		return errors.New("unexpected error. Unsupported schema type " + reflect.TypeOf(schema).Name())
	}
}
