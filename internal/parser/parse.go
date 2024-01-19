package parser

import (
	"errors"
	"gopkg.in/yaml.v3"
	"nanoc/internal/datatype"
	"nanoc/internal/symbol"
	"os"
	"strings"
)

func ParseSchema(path string) (datatype.PartialSchema, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	err = yaml.Unmarshal(b, m)
	if err != nil {
		return nil, err
	}

	var schema datatype.PartialSchema

	for k, v := range m {
		if strings.HasPrefix(k, symbol.Enum+" ") {
			s, err := parseEnumSchema(k, v)
			if err != nil {
				return nil, err
			}

			s.SchemaPath = path
			schema = s
		} else {
			body, ok := v.(map[string]interface{})
			if !ok {
				return nil, errors.New("invalid message schema body")
			}

			s, err := parseMessageSchema(k, body)
			if err != nil {
				return nil, err
			}

			s.SchemaPath = path
			schema = s
		}
	}

	return schema, nil
}

// ParseType parses the given type expression and returns the DataType describing the type and the schema that is used or nil if none is used.
// Accepts a schema map that stores user-defined types.
func ParseType(expr string, sm datatype.SchemaMap) (*datatype.DataType, datatype.Schema, error) {
	s, ok := sm[expr]
	if ok {
		t := s.DataType()
		return &t, s, nil
	}

	prim := datatype.FromIdentifier(expr)
	if prim != nil {
		return prim, nil, nil
	}

	if strings.HasSuffix(expr, symbol.Optional) {
		t, s, err := ParseType(expr[:len(expr)-len(symbol.Optional)], sm)
		if err != nil {
			return nil, nil, err
		}
		opt := datatype.NewOptionalType(t)
		return &opt, s, nil
	}

	if strings.HasSuffix(expr, symbol.Array) {
		t, s, err := ParseType(expr[:len(expr)-len(symbol.Array)], sm)
		if err != nil {
			return nil, nil, err
		}
		arr := datatype.NewArrayType(t)
		return &arr, s, nil
	}

	if strings.HasPrefix(expr, symbol.MapBracketStart) && strings.HasSuffix(expr, symbol.MapBracketEnd) {
		inner := expr[1 : len(expr)-len(symbol.MapBracketEnd)]
		ps := strings.Split(symbol.MapKeyValTypeSep, inner)
		if len(ps) != 2 {
			return nil, nil, SyntaxError{
				Msg:           "Expected a key type and a value type separated by a comma (',')",
				OffendingCode: expr,
			}
		}

		kt, _, err := ParseType(strings.TrimSpace(ps[0]), sm)
		if err != nil {
			return nil, nil, err
		}

		vt, s, err := ParseType(strings.TrimSpace(ps[1]), sm)
		if err != nil {
			return nil, nil, err
		}

		mt := datatype.NewMapType(kt, vt)
		return &mt, s, nil
	}

	return nil, nil, &SyntaxError{
		Msg:           "Invalid type expression",
		OffendingCode: expr,
	}
}
