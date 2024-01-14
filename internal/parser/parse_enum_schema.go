package parser

import (
	"errors"
	"nanoc/internal/npschema"
	"strings"
)

func parseEnumSchema(header string, body map[string]interface{}) (*npschema.PartialEnum, error) {
	ps := strings.Split(header, " ")
	if len(ps) != 2 {
		return nil, errors.New("invalid enum header declaration: expected enum EnumName, or enum EnumName::ValueType, received " + header)
	}

	ps = strings.Split(ps[1], SymbolTypeSeparator)
	if len(ps) > 2 {
		return nil, errors.New("invalid enum header declaration: expected enum EnumName, or enum EnumName::ValueType, received " + header)
	}

	schema := npschema.PartialEnum{
		Name: ps[0],
	}
	if len(ps) == 2 {
		schema.ValueTypeName = ps[1]
	}

	for k, v := range body {
		l, ok := v.(string)
		if !ok {
			return nil, errors.New("invalid enum value literal")
		}

		schema.Members = append(schema.Members, npschema.EnumMember{
			Name:         k,
			ValueLiteral: l,
		})
	}

	return &schema, nil
}
