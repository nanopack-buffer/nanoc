package parser

import "fmt"

type SyntaxError struct {
	SchemaPathAbs string
	ErrorMessage  string
}

type invalidEnumHeader struct {
	providedHeader string
}

type invalidEnumValue struct {
	memberName string
}

type invalidEnumMemberName struct {
	memberIndex int
}

type invalidMessageHeader struct {
	providedHeader string
}

type invalidTypeID struct{}

type invalidMessageFieldDeclaration struct {
	fieldName string
}

type invalidTypeDeclaration struct {
	providedDeclaration string
}

type invalidMapTypeDeclaration struct {
	providedDeclaration string
}

type mixUseOfImplicitAndExplicitFieldNumber struct{}

const errMsgInvalidMessageSchemaBody = `invalid message schema. Valid example:
MyMessage:
  my_field: string
  other_field: int32
`

func (s *SyntaxError) Error() string {
	return fmt.Sprintf("%v: syntax error: %v", s.SchemaPathAbs, s.ErrorMessage)
}

func (err *invalidEnumHeader) Error() string {
	return fmt.Sprintf("invalid enum header provided: \"%v\".", err.providedHeader)
}

func (err *invalidEnumValue) Error() string {
	return fmt.Sprintf("invalid value provided for enum member %v. It must be either a string or a number.", err.memberName)
}

func (err *invalidEnumMemberName) Error() string {
	return fmt.Sprintf("invalid enum member name at member #%d", err.memberIndex)
}

func (err *invalidMessageHeader) Error() string {
	return fmt.Sprintf("invalid message name provided: \"%v\". "+
		"Message names cannot start with a digit, contain whitespaces, or non-alphanumeric characters except underscore. "+
		"Valid message names include: \"MyMessage\", \"MyMessage::MyParentMessage\", \"My_Message\" "+
		"Invalid message names include: \"2Message4U,\", \"My Message\", \"My-Message\".", err.providedHeader)
}

func (err *invalidTypeID) Error() string {
	return "invalid type ID. Either it must be an integer or omit the type ID entirely and let me synthesize one for you."
}

func (err *invalidMessageFieldDeclaration) Error() string {
	return fmt.Sprintf("invalid message field declaration for \"%v\"", err.fieldName)
}

func (err *mixUseOfImplicitAndExplicitFieldNumber) Error() string {
	return "mix use of implicit and explicit field number. either omit field number for every field and let me decide one, or specify one for every field."
}

func (err *invalidTypeDeclaration) Error() string {
	return fmt.Sprintf("invalid type declaration: \"%v\". Refer to https://polygui.org/nanopack/data-types/ for supported types.", err.providedDeclaration)
}

func (err *invalidMapTypeDeclaration) Error() string {
	return fmt.Sprintf("invalid map type declaration: \"%v\". Map type declaration must be in the format \"<key_type:value_type>\".", err.providedDeclaration)
}
