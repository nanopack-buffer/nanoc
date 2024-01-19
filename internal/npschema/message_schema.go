package npschema

import "nanoc/internal/datatype"

// PartialMessage is a struct that contains partial information of a message schema.
// After all schemas are parsed, the resolver will be able to resolve the missing information to produce Message.
type PartialMessage struct {
	SchemaPath string

	ImportedTypeNames []string
	Name              string
	TypeID            int

	ParentMessageName string

	DeclaredFields []PartialMessageField
}

type Message struct {
	SchemaPath string

	ImportedTypes []datatype.Schema
	Name          string
	TypeID        int

	IsInherited   bool
	ChildMessages []*Message

	HasParentMessage bool
	ParentMessage    *Message

	InheritedFields []MessageField
	DeclaredFields  []MessageField
	AllFields       []MessageField
}

// PartialMessageField is used alongside PartialMessage.
// It is produced right after the parse stage when there isn't information
// to fully resolve all the information required, such as type information.
type PartialMessageField struct {
	Name     string
	TypeName string
	Number   int
}

type MessageField struct {
	Name   string
	Type   datatype.DataType
	Number int

	// The Message schema that this field is defined in.
	Schema *Message
}

func (m *Message) SchemaPathAbsolute() string {
	return m.SchemaPath
}

func (m *Message) DataType() datatype.DataType {
	return datatype.DataType{
		Identifier: m.Name,
		Kind:       datatype.Message,
		ByteSize:   datatype.DynamicSize,
		Schema:     m,
		KeyType:    nil,
		ElemType:   nil,
	}
}

func (p PartialMessage) IsPartialSchema() {}
