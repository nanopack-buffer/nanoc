package parser

import (
	"errors"
	"fmt"
	"nanoc/internal/npschema"
	"nanoc/internal/symbol"
	"strconv"
	"strings"
)

func parseEnumSchema(header string, body any) (*npschema.PartialEnum, error) {
	ps := strings.Split(header, " ")
	if len(ps) != 2 {
		return nil, errors.New("invalid enum header declaration: expected enum EnumName, or enum EnumName::ValueType, received " + header)
	}

	ps = strings.Split(ps[1], symbol.TypeSeparator)
	if len(ps) > 2 {
		return nil, errors.New("invalid enum header declaration: expected enum EnumName, or enum EnumName::ValueType, received " + header)
	}

	schema := npschema.PartialEnum{
		Name: ps[0],
	}

	if len(ps) == 2 {
		schema.ValueTypeName = ps[1]
	}

	switch body := body.(type) {
	case map[string]any:
		for k, v := range body {
			var l string
			switch v := v.(type) {
			case string:
				l = v
			case int:
				l = strconv.Itoa(v)
			default:
				return nil, SyntaxError{
					Msg:           "invalid enum value",
					OffendingCode: fmt.Sprintf("%v", v),
				}
			}

			schema.Members = append(schema.Members, npschema.EnumMember{
				Name:         k,
				ValueLiteral: l,
			})
		}

	case []any:
		for i, v := range body {
			var l string
			switch v := v.(type) {
			case string:
				l = v
			case int:
				l = strconv.Itoa(v)
			default:
				return nil, SyntaxError{
					Msg:           "invalid enum value",
					OffendingCode: fmt.Sprintf("%v", v),
				}
			}

			schema.Members = append(schema.Members, npschema.EnumMember{
				Name:         l,
				ValueLiteral: strconv.Itoa(i),
			})
		}
	}

	return &schema, nil
}
