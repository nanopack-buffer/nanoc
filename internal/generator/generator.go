package generator

import (
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
	"text/template"
)

type MessageCodeGeneratorMap map[datatype.Kind]DataTypeMessageCodeGenerator

// DataTypeMessageCodeGenerator defines the interface of a code generator that generates fragments of message code for a specific DataType.
// The fragments are put in the right places to form the complete code.
type DataTypeMessageCodeGenerator interface {
	// TypeDeclaration creates the type declaration for the given data type. The data type must be supported by this generator.
	TypeDeclaration(dataType datatype.DataType) string

	// ReadSizeExpression creates an expression to read the byte size of a variable.
	ReadSizeExpression(dataType datatype.DataType, varName string) string

	// ConstructorFieldParameter creates the constructor parameter definition of the given message field.
	ConstructorFieldParameter(field npschema.MessageField) string

	// ConstructorFieldInitializer creates the statement to initialize the given message field in the message class.
	ConstructorFieldInitializer(field npschema.MessageField) string

	// FieldDeclaration creates the field declaration for the given message field in the message class.
	FieldDeclaration(field npschema.MessageField) string

	// ReadFieldFromBuffer creates code to read the value of the given field from a NanoPack-formatted buffer.
	ReadFieldFromBuffer(field npschema.MessageField, ctx CodeContext) string

	// ReadValueFromBuffer creates code to read value to a variable from a NanoPack-formatted buffer.
	ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx CodeContext) string

	// WriteFieldToBuffer creates code to write the value of the field to a NanoPack-formatted buffer.
	WriteFieldToBuffer(field npschema.MessageField, ctx CodeContext) string

	// WriteVariableToBuffer creates code to write the value of a variable to a NanoPack-formatted buffer.
	WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx CodeContext) string

	// ToFuncMap creates a template.FuncMap for this generator that can be used in a Go template.
	ToFuncMap() template.FuncMap
}

const (
	FuncMapKeyTypeDeclaration             = "TypeDeclaration"
	FuncMapKeyReadSizeExpression          = "ReadSizeExpression"
	FuncMapKeyConstructorFieldParameter   = "ConstructorFieldParameter"
	FuncMapKeyConstructorFieldInitializer = "ConstructorFieldInitializer"
	FuncMapKeyFieldDeclaration            = "FieldDeclaration"
	FuncMapKeyReadFieldFromBuffer         = "ReadFieldFromBuffer"
	FuncMapKeyReadValueFromBuffer         = "ReadValueFromBuffer"
	FuncMapKeyWriteFieldToBuffer          = "WriteFieldToBuffer"
	FuncMapKeyWriteVariableToBuffer       = "WriteVariableToBuffer"
)
