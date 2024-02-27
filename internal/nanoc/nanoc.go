package nanoc

import (
	"errors"
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
	"nanoc/internal/parser"
	"nanoc/internal/resolver"
	"sort"
	"sync"
)

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

	schemaGenerator := schemaGeneratorMap[opts.Language]
	for _, s := range schemas {
		s := s
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := schemaGenerator(s, opts); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}()
	}

	if opts.MessageFactoryAbsFilePath != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()

			g := messageFactoryGeneratorMap[opts.Language]

			var mss []*npschema.Message
			for _, s := range schemas {
				if ms, ok := s.(*npschema.Message); ok {
					mss = append(mss, ms)
				}
			}

			sort.Slice(mss, func(i, j int) bool {
				return mss[i].TypeID < mss[j].TypeID
			})

			if err := g(mss, opts); err != nil {
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
