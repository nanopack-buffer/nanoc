package parser

import (
	"errors"
	"nanoc/internal/npschema"
	"strconv"
	"strings"
)

func parseMessageSchema(header string, body map[string]interface{}) (*npschema.PartialMessage, error) {
	ps := strings.Split(header, SymbolTypeSeparator)
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

	for k, v := range body {
		if k == SymbolTypeID {
			typeIDStr, ok := v.(string)
			if !ok {
				return nil, errors.New("non-numeric type ID received")
			}

			typeID, err := strconv.Atoi(typeIDStr)
			if err != nil {
				return nil, errors.New("invalid type ID. received: " + typeIDStr)
			}

			schema.TypeID = typeID
		} else {
			typeName, fieldNumber, err := parseFieldType(k)
			if err != nil {
				return nil, err
			}
			schema.DeclaredFields = append(schema.DeclaredFields, npschema.PartialMessageField{
				Name:     k,
				TypeName: typeName,
				Number:   fieldNumber,
			})
		}
	}

	if schema.TypeID < 0 {
		return nil, errors.New("type ID not defined")
	}

	return &schema, nil
}

func parseFieldType(str string) (string, int, error) {
	ps := strings.Split(str, SymbolFieldNumberSeparator)
	if len(ps) != 2 {
		return "", 0, errors.New("invalid field type syntax. expected field-type:field-number, received: " + str)
	}

	fieldNumber, err := strconv.Atoi(ps[1])
	if err != nil {
		return "", 0, errors.New("invalid field number. expected a valid number, received: " + ps[1])
	}

	return ps[0], fieldNumber, nil
}
