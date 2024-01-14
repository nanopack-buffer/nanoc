package npschema

import "nanoc/internal/datatype"

// PartialMessage is a struct that contains partial information of a message schema.
// After all schemas are parsed, the resolver will be able to resolve the missing information to produce Message.
type PartialMessage struct {
	ImportedTypeNames []string
	Name              string
	TypeID            int

	ParentMessageName string

	DeclaredFields []PartialMessageField
}

type Message struct {
	ImportedTypes []Schema
	Name          string
	TypeID        int

	IsInherited   bool
	ChildMessages []Message

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
}

func (m Message) isSchema() {}

func (p PartialMessage) isPartialSchema() {}
