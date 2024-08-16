package resolver

import (
	"errors"
	"fmt"
	"math"
	"nanoc/internal/datatype"
	"nanoc/internal/errs"
	"nanoc/internal/npschema"
	"nanoc/internal/parser"
	"sort"
	"strconv"
)

type ResolveResult struct {
	Schemas []datatype.Schema

	// TypesUsed stores a list of all data types used in Schemas,
	// including built-in types and user-defined types,
	TypesUsed []string
}

func Resolve(schemas []datatype.PartialSchema) (*ResolveResult, error) {
	sm := make(datatype.SchemaMap, len(schemas))
	resolvedSchemas := []datatype.Schema{}
	usedTypesMap := map[string]struct{}{}

	// create a type map from all the input schemas
	// mapping their type names to their incomplete schema
	// which will be completed later.
	for _, schema := range schemas {
		switch s := schema.(type) {
		case *npschema.PartialMessage:
			sm[s.Name] = &npschema.Message{
				SchemaPath: s.SchemaPath,
				TypeID:     s.TypeID,
				Name:       s.Name,
			}

		case *npschema.PartialEnum:
			// enum schema can be resolved directly, since an enum cannot import other types
			resolved, err := resolveEnumSchema(s, sm, usedTypesMap)
			if err != nil {
				return nil, err
			}
			resolvedSchemas = append(resolvedSchemas, resolved)
			sm[s.Name] = resolved

		case *npschema.PartialService:
			resolved, err := resolveServiceSchema(s, sm)
			if err != nil {
				return nil, errs.WrapNanocErr(err, fmt.Sprintf("service %v", s.Name))
			}
			resolvedSchemas = append(resolvedSchemas, resolved)
		}
	}

	for _, schema := range schemas {
		if pm, ok := schema.(*npschema.PartialMessage); ok {
			err := resolveMessageSchemaTypeInfo(pm, sm, usedTypesMap)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, schema := range sm {
		if ms, ok := schema.(*npschema.Message); ok {
			resolveMessageInheritance(ms)
			resolvedSchemas = append(resolvedSchemas, ms)
		}
	}

	usedTypes := make([]string, 0, len(usedTypesMap))
	for name := range usedTypesMap {
		usedTypes = append(usedTypes, name)
	}

	return &ResolveResult{
		Schemas:   resolvedSchemas,
		TypesUsed: usedTypes,
	}, nil
}

func resolveEnumSchema(partialEnum *npschema.PartialEnum, sm datatype.SchemaMap, usedTypes map[string]struct{}) (*npschema.Enum, error) {
	fullSchema := npschema.Enum{
		SchemaPath:         partialEnum.SchemaPath,
		Name:               partialEnum.Name,
		IsDefaultValueUsed: partialEnum.IsDefaultValueUsed,
	}

	fullSchema.Members = make([]npschema.EnumMember, len(partialEnum.Members))
	copy(fullSchema.Members, partialEnum.Members)

	if t := datatype.FromIdentifier(partialEnum.ValueTypeName); t != nil {
		fullSchema.ValueType = *t

		if datatype.IsInt(fullSchema.ValueType) {
			for _, m := range fullSchema.Members {
				_, err := strconv.Atoi(m.ValueLiteral)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("non-int value %v used for enum member %v in enum %v", m.ValueLiteral, m.Name, partialEnum.Name))
				}
			}
		} else if datatype.IsDouble(fullSchema.ValueType) {
			for _, m := range fullSchema.Members {
				_, err := strconv.ParseFloat(m.ValueLiteral, 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("non-double value %v used for enum member %v in enum %v", m.ValueLiteral, m.Name, partialEnum.Name))
				}
			}
		}
	} else {
		// guess type based on enum values

		nums := make([]int, 0, len(fullSchema.Members))

		ok := false
		for _, m := range partialEnum.Members {
			n, err := strconv.Atoi(m.ValueLiteral)
			if err != nil {
				fullSchema.ValueType = *datatype.FromKind(datatype.String)
				ok = true
				break
			}
			nums = append(nums, n)
		}

		if !ok {
			maxVal := -1
			for _, n := range nums {
				if n > maxVal {
					maxVal = n
				}
			}

			if maxVal <= math.MaxInt8 {
				fullSchema.ValueType = *datatype.FromKind(datatype.Int8)
			} else if maxVal <= math.MaxInt32 {
				fullSchema.ValueType = *datatype.FromKind(datatype.Int32)
			} else if maxVal <= math.MaxInt64 {
				fullSchema.ValueType = *datatype.FromKind(datatype.Int64)
			} else {
				return nil, errors.New(fmt.Sprintf("unable to determine the value type to use for enum %v", partialEnum.Name))
			}
		}
	}

	usedTypes[fullSchema.ValueType.Identifier] = struct{}{}

	return &fullSchema, nil
}

func resolveMessageSchemaTypeInfo(partialMsg *npschema.PartialMessage, sm datatype.SchemaMap, usedTypes map[string]struct{}) error {
	s, ok := sm[partialMsg.Name]
	if !ok {
		return errors.New("unexpected error when resolving " + partialMsg.Name + ": not found in type map.")
	}

	fullSchema, ok := s.(*npschema.Message)
	if !ok {
		return errors.New("unexpected error when resolving " + partialMsg.Name + ": schema type is not npschema.Message.")
	}

	imported := map[string]struct{}{}

	if partialMsg.ParentMessageName != "" {
		ps, ok := sm[partialMsg.ParentMessageName]
		if !ok {
			return errors.New("unable to resolve parent type " + partialMsg.ParentMessageName + " used in " + partialMsg.Name)
		}

		pmsg, ok := ps.(*npschema.Message)
		if !ok {
			return errors.New(partialMsg.Name + " is trying to inherit from a non-message type. This is forbidden.")
		}

		pmsg.IsInherited = true
		pmsg.ChildMessages = append(pmsg.ChildMessages, fullSchema)

		fullSchema.HasParentMessage = true
		fullSchema.ParentMessage = pmsg
		fullSchema.ImportedTypes = append(fullSchema.ImportedTypes, pmsg)

		imported[pmsg.Name] = struct{}{}
		usedTypes[pmsg.Name] = struct{}{}
	}

	for _, f := range partialMsg.DeclaredFields {
		if f.TypeName == partialMsg.Name {
			t := fullSchema.DataType()
			// self-referencing field, mark type as optional to add indirection
			fullSchema.DeclaredFields = append(fullSchema.DeclaredFields, npschema.MessageField{
				Name:   f.Name,
				Type:   datatype.NewOptionalType(t),
				Number: f.Number,
				Schema: fullSchema,
			})
			continue
		}

		t, s, err := parser.ParseType(f.TypeName, sm)
		if err != nil {
			return err
		}

		usedTypes[t.Identifier] = struct{}{}

		if s != nil {
			var name string
			switch sp := s.(type) {
			case *npschema.Message:
				name = sp.Name
			case *npschema.Enum:
				name = sp.Name
			}

			if name != partialMsg.Name {
				if _, ok := imported[name]; !ok {
					fullSchema.ImportedTypes = append(fullSchema.ImportedTypes, s)
					imported[name] = struct{}{}
				}
			}
		}

		fullSchema.DeclaredFields = append(fullSchema.DeclaredFields, npschema.MessageField{
			Name:   f.Name,
			Number: f.Number,
			Type:   *t,
			Schema: fullSchema,
		})
	}

	return nil
}

func resolveMessageInheritance(msgSchema *npschema.Message) {
	sort.Slice(msgSchema.DeclaredFields, func(i, j int) bool {
		return msgSchema.DeclaredFields[i].Number < msgSchema.DeclaredFields[j].Number
	})

	if msgSchema.HasParentMessage {
		fnum := 0
		ms := msgSchema.ParentMessage
		for ms != nil {
			// make a copy of all the fields of this schema
			// because their field numbers will be readjusted based on the inheritance chain
			// we don't want this adjustment for the original schema, because it is only specific to msgSchema
			fields := make([]npschema.MessageField, len(ms.DeclaredFields))
			copy(fields, ms.DeclaredFields)

			for i := range fields {
				fields[i].Number = fnum
				fnum += 1
			}

			msgSchema.InheritedFields = append(msgSchema.InheritedFields, fields...)

			ms = ms.ParentMessage
		}

		for i := range msgSchema.DeclaredFields {
			msgSchema.DeclaredFields[i].Number = fnum
			fnum += 1
		}

		sort.Slice(msgSchema.InheritedFields, func(i, j int) bool {
			return msgSchema.InheritedFields[i].Number < msgSchema.InheritedFields[j].Number
		})

		msgSchema.AllFields = append(msgSchema.AllFields, msgSchema.InheritedFields...)
		msgSchema.AllFields = append(msgSchema.AllFields, msgSchema.DeclaredFields...)
	} else {
		msgSchema.AllFields = make([]npschema.MessageField, len(msgSchema.DeclaredFields))
		copy(msgSchema.AllFields, msgSchema.DeclaredFields)

		sort.Slice(msgSchema.AllFields, func(i, j int) bool {
			return msgSchema.AllFields[i].Number < msgSchema.AllFields[j].Number
		})
	}

	sort.Slice(msgSchema.ChildMessages, func(i, j int) bool {
		return msgSchema.ChildMessages[i].TypeID < msgSchema.ChildMessages[j].TypeID
	})

	msgSchema.HeaderSize = len(msgSchema.AllFields)*4 + 4
}

func resolveServiceSchema(partialSchema *npschema.PartialService, sm datatype.SchemaMap) (*npschema.Service, error) {
	fullSchema := npschema.Service{
		Name:       partialSchema.Name,
		SchemaPath: partialSchema.SchemaPath,
	}

	imported := map[string]datatype.DataType{}

	for _, f := range partialSchema.DeclaredFunctions {
		fullFunc := npschema.DeclaredFunction{
			Name:               f.Name,
			ParametersByteSize: 0,
		}

		if f.ReturnTypeName != "" {
			t, s, err := parser.ParseType(f.ReturnTypeName, sm)
			if err != nil {
				return nil, errs.NewNanocError(fmt.Sprintf("Unresolved function return type %v", f.ReturnTypeName), f.Name)
			}

			if s != nil {
				imported[s.SchemaName()] = *t
			} else if t.Kind == datatype.Message {
				imported[datatype.IdentifierGenericMessage] = *t
			}

			fullFunc.ReturnType = t
		} else {
			fullFunc.ReturnType = nil
		}

		if f.ErrorTypeName != "" {
			t, s, err := parser.ParseType(f.ErrorTypeName, sm)
			if err != nil {
				return nil, errs.NewNanocError(fmt.Sprintf("Unresolved function error type %v", f.ReturnTypeName), f.Name)
			}

			if s != nil {
				imported[s.SchemaName()] = *t
			} else if t.Kind == datatype.Message {
				imported[datatype.IdentifierGenericMessage] = *t
			}

			fullFunc.ErrorType = t
		} else {
			fullFunc.ErrorType = nil
		}

		for _, param := range f.Parameters {
			t, s, err := parser.ParseType(param.TypeName, sm)
			if err != nil {
				return nil, errs.NewNanocError(fmt.Sprintf("Unresolved type %v", param.TypeName), f.Name, param.Name)
			}

			if s != nil {
				imported[s.SchemaName()] = *t
			} else if t.Kind == datatype.Message {
				imported[datatype.IdentifierGenericMessage] = *t
			}

			if t.ByteSize != datatype.DynamicSize && fullFunc.ParametersByteSize != datatype.DynamicSize {
				fullFunc.ParametersByteSize += t.ByteSize
			} else {
				fullFunc.ParametersByteSize = datatype.DynamicSize
			}

			fullFunc.Parameters = append(fullFunc.Parameters, npschema.FunctionParameter{
				Name: param.Name,
				Type: *t,
			})
		}

		fullSchema.DeclaredFunctions = append(fullSchema.DeclaredFunctions, fullFunc)
	}

	for _, t := range imported {
		fullSchema.ImportedTypes = append(fullSchema.ImportedTypes, t)
	}

	return &fullSchema, nil
}
