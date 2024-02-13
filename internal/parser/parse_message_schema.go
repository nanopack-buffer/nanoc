package parser

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"nanoc/internal/npschema"
	"nanoc/internal/symbol"
	"strconv"
	"strings"
)

func parseMessageSchema(header string, body yaml.MapSlice) (*npschema.PartialMessage, error) {
	ps := strings.Split(header, symbol.TypeSeparator)
	l := len(ps)
	if l > 2 || l <= 0 {
		return nil, errors.New("invalid message header syntax. received: " + header)
	}

	schema := npschema.PartialMessage{
		Name:   ps[0],
		TypeID: -1,
	}
	if l == 2 {
		schema.ParentMessageName = ps[1]
	}

	for i, e := range body {
		k := e.Key.(string)
		v := e.Value
		if k == symbol.TypeID {
			typeID, ok := v.(int)
			if !ok {
				return nil, SyntaxError{
					Msg:           "non-numeric type ID received",
					OffendingCode: fmt.Sprintf("%v: %v", k, v),
				}
			}
			schema.TypeID = typeID
		} else if s, ok := v.(string); ok {
			typeName, fieldNumber, err := parseFieldType(s)
			if err != nil {
				return nil, err
			}

			if fieldNumber < 0 {
				// field number is not declared, use the declaration order of the field in the schema as its field number
				// i starts with 1 because typeid field precedes message fields.
				fieldNumber = i - 1
			}
			schema.DeclaredFields = append(schema.DeclaredFields, npschema.PartialMessageField{
				Name:     k,
				TypeName: typeName,
				Number:   fieldNumber,
			})
		} else {
			return nil, SyntaxError{
				Msg:           "Invalid message schema body.",
				OffendingCode: fmt.Sprintf("%v: %v", k, v),
			}
		}
	}

	if schema.TypeID < 0 {
		return nil, errors.New("type ID not defined")
	}

	return &schema, nil
}

// parseFieldType returns the type name and the field number of the given field type declaration, e.g. string:0, int.
// field number will be returned as a negative number if it is not declared.
func parseFieldType(str string) (string, int, error) {
	ps := strings.Split(str, symbol.FieldNumberSeparator)
	if len(ps) != 2 {
		return ps[0], -1, nil
	}

	fieldNumber, err := strconv.Atoi(ps[1])
	if err != nil {
		return "", 0, errors.New("invalid field number. expected a valid number, received: " + ps[1])
	}

	return ps[0], fieldNumber, nil
}
