package parser

import (
	"errors"
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/errs"
	"nanoc/internal/symbol"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func ParseSchema(path string) (datatype.PartialSchema, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	m := yaml.MapSlice{}
	err = yaml.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	var schema datatype.PartialSchema

	for _, e := range m {
		k := e.Key.(string)
		v := e.Value
		if strings.HasPrefix(k, symbol.Enum+" ") {
			s, err := parseEnumSchema(k, v)
			if err != nil {
				return nil, err
			}

			s.SchemaPath = path
			schema = s
		} else if strings.HasPrefix(k, symbol.Service+" ") {
			body, ok := v.(yaml.MapSlice)
			if !ok {
				return nil, errors.New("invalid rpc schema body")
			}

			s, err := parseServiceSchema(k, body)
			if err != nil {
				return nil, err
			}

			s.SchemaPath = path
			schema = s
		} else {
			body, ok := v.(yaml.MapSlice)
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
		return t, s, nil
	}

	builtin := datatype.FromIdentifier(expr)
	if builtin != nil {
		return builtin, nil, nil
	}

	if strings.HasSuffix(expr, symbol.Optional) {
		t, s, err := ParseType(expr[:len(expr)-len(symbol.Optional)], sm)
		if err != nil {
			return nil, nil, errs.WrapNanocErr(err, expr)
		}
		opt := datatype.NewOptionalType(t)
		return &opt, s, nil
	}

	if strings.HasSuffix(expr, symbol.Array) {
		t, s, err := ParseType(expr[:len(expr)-len(symbol.Array)], sm)
		if err != nil {
			return nil, nil, errs.WrapNanocErr(err, expr)
		}
		arr := datatype.NewArrayType(t)
		return &arr, s, nil
	}

	if strings.HasPrefix(expr, symbol.MapBracketStart) && strings.HasSuffix(expr, symbol.MapBracketEnd) {
		inner := expr[1 : len(expr)-len(symbol.MapBracketEnd)]
		ps := strings.Split(inner, symbol.MapKeyValTypeSeperator)
		if len(ps) != 2 {
			return nil, nil, errs.NewNanocError(fmt.Sprintf("Expected a key type and a value type to be separated by '%v'", symbol.MapKeyValTypeSeperator), expr)
		}

		kt, _, err := ParseType(strings.TrimSpace(ps[0]), sm)
		if err != nil {
			return nil, nil, errs.WrapNanocErr(err, expr)
		}

		vt, s, err := ParseType(strings.TrimSpace(ps[1]), sm)
		if err != nil {
			return nil, nil, errs.WrapNanocErr(err, expr)
		}

		mt := datatype.NewMapType(kt, vt)
		return &mt, s, nil
	}

	return nil, nil, errs.NewNanocError("Invalid type expression", expr)
}
