package parser

import (
	"errors"
	"github.com/twmb/murmur3"
	"gopkg.in/yaml.v2"
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
	"nanoc/internal/symbol"
	"strconv"
	"strings"
)

const noFieldNumberProvided = -1

func parseMessageSchema(header string, body yaml.MapSlice) (*npschema.PartialMessage, error) {
	ps := strings.Split(header, symbol.TypeSeparator)
	l := len(ps)
	if l > 2 || l <= 0 {
		return nil, &invalidMessageHeader{header}
	}

	schema := npschema.PartialMessage{
		Name:   ps[0],
		TypeID: 0,
	}

	if l == 2 {
		schema.ParentMessageName = ps[1]
	}

	useImplicitFieldNumber := true
	typeIDDeclared := false

	for i, e := range body {
		k := e.Key.(string)
		v := e.Value
		if k == symbol.TypeID {
			typeID, ok := v.(int)
			if !ok {
				return nil, &invalidTypeID{}
			}
			typeIDDeclared = true
			schema.TypeID = datatype.TypeID(typeID)
		} else if s, ok := v.(string); ok {
			typeName, fieldNumber, err := parseFieldType(s)
			if err != nil {
				return nil, err
			}

			if fieldNumber == noFieldNumberProvided {
				useImplicitFieldNumber = true
				fieldNumber = i
			} else if useImplicitFieldNumber {
				return nil, &mixUseOfImplicitAndExplicitFieldNumber{}
			}

			schema.DeclaredFields = append(schema.DeclaredFields, npschema.PartialMessageField{
				Name:     k,
				TypeName: typeName,
				Number:   fieldNumber,
			})
		} else {
			return nil, &invalidMessageFieldDeclaration{
				fieldName: k,
			}
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
		return ps[0], noFieldNumberProvided, nil
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
