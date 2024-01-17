package swiftgen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"strings"
)

type mapGenerator struct {
	gm generator.MessageCodeGeneratorMap
}

func (g mapGenerator) TypeDeclaration(dataType datatype.DataType) string {
	kg := g.gm[dataType.KeyType.Kind]
	ig := g.gm[dataType.ElemType.Kind]
	return fmt.Sprintf("[%v: %v]", kg.TypeDeclaration(*dataType.KeyType), ig.TypeDeclaration(*dataType.ElemType))
}

func (g mapGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	if dataType.KeyType.ByteSize == datatype.DynamicSize || dataType.ElemType.ByteSize == datatype.DynamicSize {
		return varName + "ByteSize"
	}
	return fmt.Sprintf("%v.count * %d", varName, dataType.KeyType.ByteSize+dataType.ElemType.ByteSize)
}

func (g mapGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return fmt.Sprintf("%v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g mapGenerator) FieldInitializer(field npschema.MessageField) string {
	c := strcase.ToLowerCamel(field.Name)
	return fmt.Sprintf("self.%v = %v", c, c)
}

func (g mapGenerator) FieldDeclaration(field npschema.MessageField) string {
	return fmt.Sprintf("let %v: %v", strcase.ToLowerCamel(field.Name), g.TypeDeclaration(field.Type))
}

func (g mapGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)

	if field.Type.ElemType.ByteSize == datatype.DynamicSize || field.Type.KeyType.ByteSize == datatype.DynamicSize {
		return g.ReadValueFromBuffer(field.Type, c, ctx)
	}

	ctx.AddVariableToScope(c + "ByteSize")
	ctx.AddVariableToScope(c + "ItemCount")

	return generator.Lines(
		fmt.Sprintf("let %vByteSize = data.readSize(ofField: %d)", c, field.Number),
		fmt.Sprintf("let %vItemCount = %vByteSize / ", c, field.Type.ElemType.ByteSize+field.Type.KeyType.ByteSize),
		g.ReadValueFromBuffer(field.Type, c, ctx))
}

func (g mapGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	kg := g.gm[dataType.KeyType.Kind]
	ig := g.gm[dataType.ElemType.Kind]

	var l0 string
	if !ctx.IsVariableInScope(varName + "ItemCount") {
		l0 = generator.Lines(
			fmt.Sprintf("let %vItemCount = data.readSize(at: ptr)"),
			"ptr += 4")
	}

	var l1 string
	if !ctx.IsVariableInScope(varName) {
		l1 = fmt.Sprintf("var %v: %v = [:]", varName, g.TypeDeclaration(dataType))
	}

	lv := ctx.NewLoopVar()
	kv := lv + "Key"
	iv := lv + "Value"
	ctx.AddVariableToScope(kv)
	ctx.AddVariableToScope(iv)

	ls := generator.Lines(
		l0,
		l1,
		fmt.Sprintf("%v.reserveCapacity(%vItemCount)", varName),
		fmt.Sprintf("for %v in 0..<%vItemCount {", lv),
		kg.ReadValueFromBuffer(*dataType.KeyType, kv, ctx),
		ig.ReadValueFromBuffer(*dataType.ElemType, iv, ctx),
		fmt.Sprintf("%v[%v] = %v", varName, kv, iv),
		"}")

	ctx.RemoveVariableFromScope(lv)
	ctx.RemoveVariableFromScope(kv)
	ctx.RemoveVariableFromScope(iv)

	return ls
}

func (g mapGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	c := strcase.ToLowerCamel(field.Name)
	kg := g.gm[field.Type.KeyType.Kind]
	ig := g.gm[field.Type.ElemType.Kind]

	if field.Type.ElemType.ByteSize == datatype.DynamicSize || field.Type.KeyType.ByteSize == datatype.DynamicSize {
		return generator.Lines(
			g.WriteVariableToBuffer(field.Type, c, ctx),
			fmt.Sprintf("data.write(size: %vByteSize, ofField: %d)", c, field.Number))
	}

	kv := ctx.NewLoopVarWithSuffix("Key")
	iv := ctx.NewLoopVarWithSuffix("Value")

	ls := generator.Lines(
		fmt.Sprintf("data.write(size: %v.count * %d, ofField: %d)", c, field.Type.ElemType.ByteSize+field.Type.KeyType.ByteSize, field.Number),
		fmt.Sprintf("for (%v, %v) in %v {", kv, iv, c),
		kg.WriteVariableToBuffer(*field.Type.KeyType, kv, ctx),
		ig.WriteVariableToBuffer(*field.Type.ElemType, iv, ctx),
		"}")

	ctx.RemoveVariableFromScope(kv)
	ctx.RemoveVariableFromScope(iv)

	return ls
}

func (g mapGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	kg := g.gm[dataType.KeyType.Kind]
	ig := g.gm[dataType.ElemType.Kind]

	var l1 string
	var l5 string
	var l6 string
	if dataType.KeyType.ByteSize == datatype.DynamicSize || dataType.ElemType.ByteSize == datatype.DynamicSize {
		b := strings.Builder{}
		b.WriteString(fmt.Sprintf("let %vByteSize = 4", varName))
		if dataType.KeyType.ByteSize != datatype.DynamicSize {
			b.WriteString(fmt.Sprintf(" + %d", dataType.KeyType.ByteSize))
			l5 = fmt.Sprintf("%vByteSize += %d", dataType.KeyType.ByteSize)
		}
		if dataType.ElemType.ByteSize != datatype.DynamicSize {
			b.WriteString(fmt.Sprintf(" + %d", dataType.ElemType.ByteSize))
			l6 = fmt.Sprintf("%vByteSize += %d", dataType.ElemType.ByteSize)
		}
		l1 = b.String()
	}

	kv := ctx.NewLoopVarWithSuffix("Key")
	iv := ctx.NewLoopVarWithSuffix("Value")

	ls := generator.Lines(
		fmt.Sprintf("data.append(size: %v.count)", varName),
		l1,
		fmt.Sprintf("for (%v, %v) in %v {", kv, iv, varName),
		kg.WriteVariableToBuffer(*dataType.KeyType, kv, ctx),
		ig.WriteVariableToBuffer(*dataType.ElemType, iv, ctx),
		l5,
		l6,
		"}")

	ctx.RemoveVariableFromScope(kv)
	ctx.RemoveVariableFromScope(iv)

	return ls
}
