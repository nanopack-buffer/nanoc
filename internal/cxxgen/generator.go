package cxxgen

import (
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
)

type cxxCodeFragmentGeneratorMap map[datatype.Kind]cxxCodeFragmentGenerator

// cxxCodeFragmentGenerator defines an interface for a generator that generates code fragments
// for a particular data type that will be added into the final generated code.
type cxxCodeFragmentGenerator interface {
	// TypeDeclaration creates the type declaration for the given data type. The data type must be supported by this generator.
	TypeDeclaration(dataType datatype.DataType) string

	ParameterDeclaration(dataType datatype.DataType, paramName string) string

	RValue(dataType datatype.DataType, argName string) string

	// ReadSizeExpression creates an expression to read the byte size of a variable.
	ReadSizeExpression(dataType datatype.DataType, varName string) string

	// ConstructorFieldParameter creates the constructor parameter definition of the given message field.
	ConstructorFieldParameter(field npschema.MessageField) string

	// FieldInitializer creates the statement to initialize the given message field in the message class.
	FieldInitializer(field npschema.MessageField) string

	// FieldDeclaration creates the field declaration for the given message field in the message class.
	FieldDeclaration(field npschema.MessageField) string

	// ReadFieldFromBuffer creates code to read the value of the given field from a NanoPack-formatted buffer.
	ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string

	// ReadValueFromBuffer creates code to read value to a variable from a NanoPack-formatted buffer.
	ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string

	// WriteFieldToBuffer creates code to write the value of the field to a NanoPack-formatted buffer.
	WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string

	// WriteVariableToBuffer creates code to write the value of a variable to a NanoPack-formatted buffer.
	WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string
}
