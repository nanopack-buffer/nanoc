package parser

import (
	"errors"
	"gopkg.in/yaml.v3"
	"nanoc/internal/npschema"
	"os"
	"strings"
)

func Parse(path string) (npschema.PartialSchema, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	err = yaml.Unmarshal(b, m)
	if err != nil {
		return nil, err
	}

	var schema npschema.PartialSchema

	for k, v := range m {
		if strings.HasPrefix(SymbolEnum+" ", k) {
			body, ok := v.(map[string]interface{})
			if !ok {
				return nil, errors.New("invalid enum schema body")
			}

			schema, err = parseEnumSchema(k, body)
			if err != nil {
				return nil, err
			}
		} else {
			body, ok := v.(map[string]interface{})
			if !ok {
				return nil, errors.New("invalid message schema body")
			}

			schema, err = parseMessageSchema(k, body)
			if err != nil {
				return nil, err
			}
		}
	}

	return schema, nil
}
