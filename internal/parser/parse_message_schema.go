package parser

import (
	"errors"
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/errs"
	"nanoc/internal/npschema"
	"nanoc/internal/symbol"
	"strconv"
	"strings"

	"github.com/twmb/murmur3"
	"gopkg.in/yaml.v2"
)

func parseMessageSchema(header string, body yaml.MapSlice) (*npschema.PartialMessage, error) {
	ps := strings.Split(header, symbol.TypeSeparator)
	l := len(ps)
	if l > 2 || l <= 0 {
		return nil, errs.NewNanocError("Invalid message header syntax", header)
	}

	schema := npschema.PartialMessage{
		Name:   ps[0],
		TypeID: 0,
	}
	typeIDDeclared := false

	if l == 2 {
		schema.ParentMessageName = ps[1]
	}

	for i, e := range body {
		k := e.Key.(string)
		v := e.Value
		if k == symbol.TypeID {
			typeID, ok := v.(int)
			if !ok {
				return nil, errs.NewNanocError("non-numeric type ID received", fmt.Sprintf("%v message", schema.Name))
			}
			typeIDDeclared = true
			schema.TypeID = datatype.TypeID(typeID)
		} else if s, ok := v.(string); ok {
			schema.DeclaredFields = append(schema.DeclaredFields, npschema.PartialMessageField{
				Name:     k,
				TypeName: s,
				Number:   i,
			})
		} else {
			return nil, errs.NewNanocError("Invalid message schema body", fmt.Sprintf("%v message", schema.Name))
		}
	}

	if !typeIDDeclared {
		schema.TypeID = calculateTypeID(schema.Name)
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

func calculateTypeID(msgName string) datatype.TypeID {
	return murmur3.Sum32([]byte(msgName))
}
