package swiftgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"

	"github.com/iancoleman/strcase"
)

type messageGenerator struct{}

func (g messageGenerator) TypeDeclaration(dataType datatype.DataType) string {
	if dataType.Schema == nil {
		return "NanoPackMessage"
	}
	return dataType.Identifier
}

func (g messageGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	return fmt.Sprintf("%vByteSize", varName)
}

func (g messageGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("%v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g messageGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("self.%v = %v", c, c)
}

func (g messageGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("let %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g messageGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	var v string
	var l4 string
	if ctx.IsVariableInScope(c) {
		v = c + "_"
		l4 = fmt.Sprintf("%v = %v", c, v)
	} else {
		v = c
	}

	var ctor string
	if field.Type.Schema == nil {
		ctor = fmt.Sprintf("makeNanoPackMessage(from: data[ptr...])")
	} else if field.Type.Schema.(*npschema.Message).IsInherited {
		ctor = fmt.Sprintf("%v.from(data: data[ptr...])", g.TypeDeclaration(field.Type))
	} else {
		ctor = fmt.Sprintf("%v(data: data[ptr...])", g.TypeDeclaration(field.Type))
	}

	return generator.Lines(
		fmt.Sprintf("let %vByteSize = data.readSize(ofField: %d)", c, field.Number),
		fmt.Sprintf("guard let %v = %v else {", v, ctor),
		"    return nil",
		"}",
		l4,
		fmt.Sprintf("ptr += %vByteSize", c))
}

func (g messageGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	var v string
	var l4 string
	if ctx.IsVariableInScope(varName) {
		v = varName + "_"
		l4 = fmt.Sprintf("%v = %v", varName, v)
	} else {
		v = varName
	}

	var ctor string
	if dataType.Schema == nil {
		ctor = fmt.Sprintf("makeNanoPackMessage(from: data[ptr...])")
	} else if dataType.Schema.(*npschema.Message).IsInherited {
		ctor = fmt.Sprintf("%v.from(data: data[ptr...], bytesRead: &%vByteSize)", g.TypeDeclaration(dataType), varName)
	} else {
		ctor = fmt.Sprintf("%v(data: data[ptr...], bytesRead: &%vByteSize)", g.TypeDeclaration(dataType), varName)
	}

	return generator.Lines(
		fmt.Sprintf("var %vByteSize = 0", varName),
		fmt.Sprintf("guard let %v = %v(data: data[ptr...], bytesRead: &%vByteSize) else {", v, ctor, varName),
		"    return nil",
		"}",
		l4,
		fmt.Sprintf("ptr += %vByteSize", varName))
}

func (g messageGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	lc := strcase.ToLowerCamel(field.Name)
	return generator.Lines(
		fmt.Sprintf("let %vByteSize = %v.write(to: &data, offset: data.count)", lc, lc),
		fmt.Sprintf("data.write(size: %vByteSize, ofField: %d, offset: offset)", lc, field.Number))
}

func (g messageGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	return fmt.Sprintf("let %vByteSize = %v.write(to: &data, offset: data.count)", varName, varName)
}
