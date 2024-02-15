package swiftgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"strings"
)

type arrayGenerator struct {
	gm generator.MessageCodeGeneratorMap
}

func (g arrayGenerator) TypeDeclaration(dataType datatype.DataType) string {
	ig := g.gm[dataType.ElemType.Kind]
	return fmt.Sprintf("[%v]", ig.TypeDeclaration(*dataType.ElemType))
}

func (g arrayGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	if dataType.ElemType.ByteSize == datatype.DynamicSize {
		return varName + "ByteSize"
	}
	return fmt.Sprintf("%v.count * %d", varName, dataType.ElemType.ByteSize)
}

func (g arrayGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("%v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g arrayGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("self.%v = %v", c, c)
}

func (g arrayGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("let %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g arrayGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)

	if field.Type.ElemType.ByteSize != datatype.DynamicSize {
		ctx.AddVariableToScope(c + "ByteSize")
		ctx.AddVariableToScope(c + "ItemCount")
		return generator.Lines(
			fmt.Sprintf("let %vByteSize = data.readSize(ofField: %d)", c, field.Number),
			fmt.Sprintf("let %vItemCount = %vSize / %d", c, c, field.Type.ElemType.ByteSize),
			g.ReadValueFromBuffer(field.Type, c, ctx))
	}

	return g.ReadValueFromBuffer(field.Type, c, ctx)
}

func (g arrayGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.ElemType.Kind]

	if dataType.ElemType.ByteSize == datatype.DynamicSize {
		var l2 string
		if !ctx.IsVariableInScope(varName) {
			l2 = fmt.Sprintf("var %v: %v = []", varName, g.TypeDeclaration(dataType))
		}

		lv := ctx.NewLoopVar()

		ls := generator.Lines(
			fmt.Sprintf("let %vItemCount = data.readSize(at: ptr)", varName),
			"ptr += 4",
			l2,
			fmt.Sprintf("%v.reserveCapacity(%vItemCount)", varName, varName),
			fmt.Sprintf("for _ in 0..<%vItemCount {", varName),
			ig.ReadValueFromBuffer(*dataType.ElemType, lv+"Item", ctx),
			fmt.Sprintf("    %v.append(%v)", varName, lv+"Item"),
			"}")

		ctx.RemoveVariableFromScope(lv)

		return ls
	}

	var l0 string
	if !ctx.IsVariableInScope(varName + "ByteSize") {
		l0 = generator.Lines(
			fmt.Sprintf("let %vItemCount = data.readSize(at: ptr)", varName),
			"ptr += 4",
			fmt.Sprintf("let %vByteSize = %vItemCount * %d", varName, varName, dataType.ElemType.ByteSize))
	}

	// code for binding memory to type
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("    %v($0.bindMemory(to: %v.self)", g.TypeDeclaration(dataType), ig.TypeDeclaration(*dataType.ElemType)))
	if datatype.IsNumber(*dataType.ElemType) {
		b.WriteString(".lazy.map{ $0.littleEndian }")
	}
	b.WriteString(")")

	if ctx.IsVariableInScope(varName) {
		return generator.Lines(
			l0,
			fmt.Sprintf("%v = data[ptr..<ptr + %vByteSize].withUnsafeBytes {", varName, varName),
			b.String(),
			"}",
			fmt.Sprintf("ptr += %vByteSize", varName))
	}

	return generator.Lines(
		l0,
		fmt.Sprintf("let %v = data[ptr..<ptr + %vByteSize].withUnsafeBytes {", varName, varName),
		b.String(),
		"}",
		fmt.Sprintf("ptr += %vByteSize", varName))
}

func (g arrayGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	ig := g.gm[field.Type.ElemType.Kind]
	c := strcase.ToLowerCamel(field.Name)

	if field.Type.ElemType.ByteSize == datatype.DynamicSize {
		return generator.Lines(
			g.WriteVariableToBuffer(field.Type, c, ctx),
			fmt.Sprintf("data.write(size: %vByteSize, ofField: %d, offset: offset)", c, field.Number))
	}

	lv := ctx.NewLoopVar()

	ls := generator.Lines(
		fmt.Sprintf("data.write(size: %v.count * %d, ofField: %d, offset: offset)", c, field.Type.ElemType.ByteSize, field.Number),
		fmt.Sprintf("for %v in %v {", lv, c),
		ig.WriteVariableToBuffer(*field.Type.ElemType, lv, ctx),
		"}")

	ctx.RemoveVariableFromScope(lv)

	return ls
}

func (g arrayGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.ElemType.Kind]

	if dataType.ElemType.ByteSize == datatype.DynamicSize {
		lv := ctx.NewLoopVar()
		ls := generator.Lines(
			fmt.Sprintf("data.append(size: %v.count)", varName),
			fmt.Sprintf("var %vByteSize: Size = 4", varName),
			fmt.Sprintf("for %v in %v {", lv, varName),
			ig.WriteVariableToBuffer(dataType, lv, ctx),
			fmt.Sprintf("%vByteSize += %v", varName, ig.ReadSizeExpression(*dataType.ElemType, lv)),
			"}")
		ctx.RemoveVariableFromScope(lv)
		return ls
	}

	lv := ctx.NewLoopVar()
	ls := generator.Lines(
		fmt.Sprintf("data.append(size: %v.count)", varName),
		fmt.Sprintf("for %v in %v {", lv, varName),
		ig.WriteVariableToBuffer(dataType, lv, ctx),
		"}")
	ctx.RemoveVariableFromScope(lv)
	return ls
}
