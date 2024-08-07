package tsgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"

	"github.com/iancoleman/strcase"
)

type mapGenerator struct {
	gm generator.MessageCodeGeneratorMap
}

func (g mapGenerator) TypeDeclaration(dataType datatype.DataType) string {
	kg := g.gm[dataType.KeyType.Kind]
	ig := g.gm[dataType.ElemType.Kind]
	return fmt.Sprintf("Map<%v, %v>", kg.TypeDeclaration(*dataType.KeyType), ig.TypeDeclaration(*dataType.ElemType))
}

func (g mapGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	if dataType.ElemType.ByteSize == datatype.DynamicSize || dataType.KeyType.ByteSize == datatype.DynamicSize {
		return varName + "ByteLength"
	}
	return fmt.Sprintf("%v.size * %d", varName, dataType.ElemType.ByteSize+dataType.KeyType.ByteSize)
}

func (g mapGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g mapGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("this.%v = %v;", c, c)
}

func (g mapGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("public %v: %v;", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g mapGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	if field.Type.ElemType.ByteSize == datatype.DynamicSize || field.Type.KeyType.ByteSize == datatype.DynamicSize {
		return g.ReadValueFromBuffer(field.Type, c, ctx)
	}
	return generator.Lines(
		fmt.Sprintf("const %vByteLength = reader.readFieldSize(%d, offset)", c, field.Number),
		fmt.Sprintf("const %vItemCount = %vByteLength / %d", c, c, field.Type.ElemType.ByteSize+field.Type.KeyType.ByteSize),
		g.ReadValueFromBuffer(field.Type, c, ctx))
}

func (g mapGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	i32g := g.gm[datatype.Int32]
	kg := g.gm[dataType.KeyType.Kind]
	ig := g.gm[dataType.ElemType.Kind]

	var l0 string
	if ctx.IsVariableInScope(varName) {
		l0 = fmt.Sprintf("%v = new %v();", varName, g.TypeDeclaration(dataType))
	} else {
		l0 = fmt.Sprintf("const %v = new %v();", varName, g.TypeDeclaration(dataType))
	}

	var l1 string
	if !ctx.IsVariableInScope(varName + "ItemCount") {
		l1 = i32g.ReadValueFromBuffer(*datatype.FromKind(datatype.Int32), varName+"ItemCount", ctx)
	}

	lv := ctx.NewLoopVar()

	ls := generator.Lines(
		l0,
		l1,
		fmt.Sprintf("for (let %v = 0; %v < %vItemCount; %v++) {", lv, lv, varName, lv),
		kg.ReadValueFromBuffer(*dataType.KeyType, lv+"Key", ctx),
		ig.ReadValueFromBuffer(*dataType.ElemType, lv+"Item", ctx),
		fmt.Sprintf("%v.set(%vKey, %vItem);", varName, lv, lv),
		"}")

	return ls
}

func (g mapGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	kg := g.gm[field.Type.KeyType.Kind]
	ig := g.gm[field.Type.ElemType.Kind]

	kv := ctx.NewLoopVarWithSuffix("Key")
	iv := ctx.NewLoopVarWithSuffix("Item")

	var lines string
	if field.Type.ElemType.ByteSize == datatype.DynamicSize || field.Type.KeyType.ByteSize == datatype.DynamicSize {
		lines = generator.Lines(
			fmt.Sprintf("writer.appendInt32(this.%v.size);", c),
			fmt.Sprintf("let %vByteLength = 4;", c),
			fmt.Sprintf("this.%v.forEach((%v, %v) => {", c, iv, kv),
			kg.WriteVariableToBuffer(*field.Type.KeyType, kv, ctx),
			ig.WriteVariableToBuffer(*field.Type.ElemType, iv, ctx),
			fmt.Sprintf("%vByteLength += (%v + %v);", c, kg.ReadSizeExpression(*field.Type.KeyType, kv), ig.ReadSizeExpression(*field.Type.ElemType, iv)),
			"});",
		)
	} else {
		lines = generator.Lines(
			fmt.Sprintf("const %vByteLength = this.%v.size * %d", c, c, field.Type.ElemType.ByteSize+field.Type.KeyType.ByteSize),
			fmt.Sprintf("writer.writeFieldSize(%d, %vByteLength, offset);", field.Number, c),
			fmt.Sprintf("this.%v.forEach((%v, %v) => {", c, kv, iv),
			kg.WriteVariableToBuffer(*field.Type.KeyType, kv, ctx),
			ig.WriteVariableToBuffer(*field.Type.ElemType, iv, ctx),
			"});",
			fmt.Sprintf("bytesWritten += %vByteLength", c))
	}

	ctx.RemoveVariableFromScope(kv)
	ctx.RemoveVariableFromScope(iv)

	return lines
}

func (g mapGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	kg := g.gm[dataType.KeyType.Kind]
	ig := g.gm[dataType.ElemType.Kind]

	kv := ctx.NewLoopVarWithSuffix("Key")
	iv := ctx.NewLoopVarWithSuffix("Item")

	var byteLengthCounter string
	var addByteLength string
	if dataType.KeyType.ByteSize == datatype.DynamicSize || dataType.ElemType.ByteSize == datatype.DynamicSize {
		byteLengthCounter = fmt.Sprintf("let %vByteLength = 4;", varName)
		addByteLength =
			fmt.Sprintf("%vByteLength += (%v + %v);", varName, kg.ReadSizeExpression(*dataType.KeyType, kv), ig.ReadSizeExpression(*dataType.ElemType, iv))
	}

	ls := generator.Lines(
		fmt.Sprintf("writer.appendInt32(%v.size);", varName),
		byteLengthCounter,
		fmt.Sprintf("%v.forEach((%v, %v) => {", varName, kv, iv),
		kg.WriteVariableToBuffer(*dataType.KeyType, varName, ctx),
		ig.WriteVariableToBuffer(*dataType.ElemType, iv, ctx),
		addByteLength,
		"});",
	)

	ctx.RemoveVariableFromScope(kv)
	ctx.RemoveVariableFromScope(iv)

	return ls
}
