package resolver

import (
	"container/list"
	"errors"
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
)

type typeMap map[string]npschema.Schema

func Resolve(schemas []npschema.PartialSchema) ([]npschema.Schema, error) {
	tm := make(typeMap, len(schemas))

	// create a type map from all the input schemas
	// mapping their type names to their incomplete schema
	// which will be completed later.
	for _, schema := range schemas {
		switch s := schema.(type) {
		case *npschema.PartialMessage:
			tm[s.Name] = &npschema.Message{}

		case *npschema.PartialEnum:
			tm[s.Name] = &npschema.Enum{}
			// enum schema can be resolved directly, since an enum cannot import other types
			err := resolveEnumSchema(s, tm)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, schema := range schemas {
		if pm, ok := schema.(*npschema.PartialMessage); ok {
			err := resolveMessageSchemaTypeInfo(pm, tm)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, schema := range tm {
		if ms, ok := schema.(*npschema.Message); ok {
			resolveMessageInheritance(ms, tm)
		}
	}

	ss := make([]npschema.Schema, len(tm))
	for _, sp := range tm {
		switch s := sp.(type) {
		case *npschema.Message:
			ss = append(ss, *s)
		case *npschema.Enum:
			ss = append(ss, *s)
		}
	}

	return ss, nil
}

func resolveEnumSchema(partialEnum *npschema.PartialEnum, tm typeMap) error {
	s, ok := tm[partialEnum.Name]
	if !ok {
		return errors.New("unexpected error when resolving " + partialEnum.Name + ": not found in type map.")
	}

	fullSchema, ok := s.(*npschema.Enum)
	if !ok {
		return errors.New("unexpected error when resolving " + partialEnum.Name + ": schema type is not npschema.Enum.")
	}

	if partialEnum.ValueTypeName == "" {
		// empty string means that the value type of the enum is not declared
		// when the value type is not declared, the enum will be implicitly backed by an int8.
		fullSchema.ValueType = datatype.Int8
	} else if t := datatype.FromIdentifier(partialEnum.ValueTypeName); t != nil {
		fullSchema.ValueType = *t
	} else {
		return errors.New("unsupported type " + partialEnum.ValueTypeName + " used for values of " + partialEnum.Name)
	}

	fullSchema.Name = partialEnum.Name
	copy(fullSchema.Members, partialEnum.Members)

	return nil
}

func resolveMessageSchemaTypeInfo(partialMsg *npschema.PartialMessage, tm typeMap) error {
	s, ok := tm[partialMsg.Name]
	if !ok {
		return errors.New("unexpected error when resolving " + partialMsg.Name + ": not found in type map.")
	}

	fullSchema, ok := s.(*npschema.Message)
	if !ok {
		return errors.New("unexpected error when resolving " + partialMsg.Name + ": schema type is not npschema.Message.")
	}

	if partialMsg.ParentMessageName != "" {
		ps, ok := tm[partialMsg.ParentMessageName]
		if !ok {
			return errors.New("unable to resolve parent type " + partialMsg.ParentMessageName + " used in " + partialMsg.Name)
		}

		pmsg, ok := ps.(*npschema.Message)
		if !ok {
			return errors.New(partialMsg.Name + " is trying to inherit from a non-message type. This is forbidden.")
		}

		fullSchema.HasParentMessage = true
		fullSchema.ParentMessage = pmsg
	}

	for _, f := range partialMsg.DeclaredFields {
		if f.TypeName == partialMsg.Name {
			// self-referencing field
			fullSchema.DeclaredFields = append(fullSchema.DeclaredFields, npschema.MessageField{
				Name:   f.TypeName,
				Type:   *datatype.FromSchema(fullSchema),
				Number: f.Number,
			})
			continue
		}

		t := datatype.FromIdentifier(f.TypeName)
		imported := map[string]bool{}
		if t == nil {
			s, ok := tm[f.TypeName]
			if !ok {
				return errors.New("unable to resolve type " + f.TypeName + " for field " + f.Name + " declared in " + partialMsg.Name)
			}

			if _, ok := imported[f.TypeName]; !ok {
				fullSchema.ImportedTypes = append(fullSchema.ImportedTypes, s)
				imported[f.TypeName] = true
			}

			fullSchema.DeclaredFields = append(fullSchema.DeclaredFields, npschema.MessageField{
				Name:   f.Name,
				Number: f.Number,
				Type:   *datatype.FromSchema(s),
			})
		} else {
			fullSchema.DeclaredFields = append(fullSchema.DeclaredFields, npschema.MessageField{
				Name:   f.Name,
				Type:   *t,
				Number: f.Number,
			})
		}
	}

	return nil
}

func resolveMessageInheritance(msgSchema *npschema.Message, tm typeMap) {
	if msgSchema.HasParentMessage {
		ic := list.New()
		cur := msgSchema
		for cur != nil {
			ic.PushFront(cur)
			cur = cur.ParentMessage
		}

		fnum := 0
		n := ic.Front()
		for n != nil {
			ms := n.Value.(*npschema.Message)
			// make a copy of all the fields of this schema
			// because their field numbers will be readjusted based on the inheritance chain
			// we don't want this adjustment for the original schema, because it is only specific to msgSchema
			fields := make([]npschema.MessageField, len(ms.DeclaredFields))
			copy(fields, ms.DeclaredFields)

			for i := range fields {
				fields[i].Number = fnum
				fnum += 1
			}

			next := n.Next()
			if next != nil {
				msgSchema.InheritedFields = append(msgSchema.InheritedFields, fields...)
			}

			n = next
		}

		msgSchema.AllFields = append(msgSchema.AllFields, msgSchema.InheritedFields...)
		msgSchema.AllFields = append(msgSchema.AllFields, msgSchema.DeclaredFields...)
	} else {
		copy(msgSchema.AllFields, msgSchema.InheritedFields)
	}
}
