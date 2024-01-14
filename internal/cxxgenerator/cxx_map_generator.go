package cxxgenerator

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"strings"
	"text/template"
)

type CxxMapGenerator struct {
	gm generator.MessageCodeGeneratorMap
}

func (g CxxMapGenerator) TypeDeclaration(dataType datatype.DataType) string {
	kg := g.gm[dataType.KeyType.Kind]
	ig := g.gm[dataType.ElemType.Kind]
	return fmt.Sprintf("std::unordered_map<%v, %v>", kg.TypeDeclaration(*dataType.KeyType), ig.TypeDeclaration(*dataType.ElemType))
}

func (g CxxMapGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	if dataType.KeyType.ByteSize != datatype.DynamicSize && dataType.ElemType.ByteSize != datatype.DynamicSize {
		return fmt.Sprintf("%v.size() * %d", varName, dataType.KeyType.ByteSize+dataType.ElemType.ByteSize)
	}
	return varName + "_byte_size"
}

func (g CxxMapGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name)
}

func (g CxxMapGenerator) ConstructorFieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	return fmt.Sprintf("%v(std::move(%v))", s, s)
}

func (g CxxMapGenerator) FieldDeclaration(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name) + ";"
}

func (g CxxMapGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)
	var l1 string
	if field.Type.ElemType.ByteSize != datatype.DynamicSize {
		// for arrays with fixed size items, the number of elements in the array can be calculated.
		l1 = fmt.Sprintf("const int32_t %v_map_size = %v_byte_size / %d;", s, s, field.Type.KeyType.ByteSize+field.Type.ElemType.ByteSize)
	}
	return generator.Lines(
		fmt.Sprintf("const int32_t %v_byte_size = reader.read_field_size(%d);", s, field.Number),
		l1,
		g.ReadValueFromBuffer(field.Type, s, ctx),
		fmt.Sprintf("this->%v = %v;", s, s))
}

func (g CxxMapGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	i32g := g.gm[datatype.Int32]
	kg := g.gm[dataType.KeyType.Kind]
	ig := g.gm[dataType.ElemType.Kind]
	mapSizeVar := varName + "_map_size"

	// If the number of elements in the vector is not read previously,
	// generate code to read it here.
	var l0 string
	if !ctx.IsVariableInScope(mapSizeVar) {
		l0 = i32g.ReadValueFromBuffer(*datatype.FromKind(datatype.Int32), mapSizeVar, ctx)
	}

	var l1 string
	if ctx.IsVariableInScope(varName) {
		l1 = fmt.Sprintf("%v = %v()", varName, g.TypeDeclaration(dataType))
	} else {
		l1 = fmt.Sprintf("%v %v()", varName, g.TypeDeclaration(dataType))
	}

	lv := ctx.NewLoopVar()
	ls := generator.Lines(
		l0,
		l1,
		fmt.Sprintf("%v.reserve(%v);", varName, mapSizeVar),
		fmt.Sprintf("for (int %v = 0; %v < %v; %v++) {", lv, lv, mapSizeVar, lv),
		kg.ReadValueFromBuffer(*dataType.KeyType, lv+"_key", ctx),
		ig.ReadValueFromBuffer(*dataType.ElemType, lv+"_value", ctx),
		fmt.Sprintf("%v.insert({%v_key, %v_value});", varName, lv, lv),
		"}")

	ctx.RemoveVariableFromScope(lv)

	return ls
}

func (g CxxMapGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)

	if field.Type.ElemType.ByteSize != datatype.DynamicSize && field.Type.KeyType.ByteSize != datatype.DynamicSize {
		lv := ctx.NewLoopVar()
		lvk := lv + "_key"
		lvv := lv + "_value"
		kg := g.gm[field.Type.KeyType.Kind]
		ig := g.gm[field.Type.ElemType.Kind]

		ctx.AddVariableToScope(lvk)
		ctx.AddVariableToScope(lvv)

		ls := generator.Lines(
			fmt.Sprintf("writer.write_field_size(%d, %v.size() * %d;", field.Number, s, field.Type.ElemType.ByteSize+field.Type.KeyType.ByteSize),
			fmt.Sprintf("for (const auto &%v : %v) {", lv, s),
			fmt.Sprintf("  const auto %v = %v.first;", lvk, lv),
			fmt.Sprintf("  const auto %v = %v.second;", lvv, lv),
			kg.WriteVariableToBuffer(*field.Type.KeyType, lvk, ctx),
			ig.WriteVariableToBuffer(*field.Type.ElemType, lvv, ctx),
			"}")

		ctx.RemoveVariableFromScope(lv)
		ctx.RemoveVariableFromScope(lvk)
		ctx.RemoveVariableFromScope(lvv)

		return ls
	}

	return generator.Lines(
		g.WriteVariableToBuffer(field.Type, s, ctx),
		fmt.Sprintf("writer.write_field_size(%d, %v);", field.Number, s+"_byte_size"))
}

func (g CxxMapGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	kg := g.gm[dataType.KeyType.Kind]
	ig := g.gm[dataType.ElemType.Kind]
	i32g := g.gm[datatype.Int32]
	mapSizeVar := varName + "_map_size"
	isDynamicSize := dataType.ElemType.ByteSize == datatype.DynamicSize || dataType.KeyType.ByteSize == datatype.DynamicSize

	l0 := fmt.Sprintf("const size_t %v = %v.size();", mapSizeVar, varName)
	l1 := i32g.WriteVariableToBuffer(*datatype.FromKind(datatype.Int32), mapSizeVar, ctx)

	var l2 string
	if isDynamicSize {
		b := strings.Builder{}

		b.WriteString(fmt.Sprintf("int32_t %v_byte_size = 4", varName))
		if dataType.KeyType.ByteSize == datatype.DynamicSize {
			b.WriteString(fmt.Sprintf(" + %v * %d", mapSizeVar, dataType.KeyType.ByteSize))
		}
		if dataType.ElemType.ByteSize == datatype.DynamicSize {
			b.WriteString(fmt.Sprintf(" + %v * %d", mapSizeVar, dataType.ElemType.ByteSize))
		}
		b.WriteString(";")

		l2 = b.String()
	}

	lv := ctx.NewLoopVar()
	lvk := lv + "_key"
	lvv := lv + "_value"

	ctx.AddVariableToScope(lvk)
	ctx.AddVariableToScope(lvv)

	ls := generator.Lines(
		l0,
		l1,
		l2,
		fmt.Sprintf("for (const auto &%v : %v) {", lv, varName),
		fmt.Sprintf("auto %v = %v.first", lvk, lv),
		fmt.Sprintf("auto %v = %v.second", lvv, lv),
		kg.WriteVariableToBuffer(*dataType.KeyType, lvk, ctx),
		ig.WriteVariableToBuffer(*dataType.ElemType, lvv, ctx))

	var l8 string
	if dataType.KeyType.ByteSize == datatype.DynamicSize {
		l8 = fmt.Sprintf("%v_byte_size += %v;", varName, kg.ReadSizeExpression(*dataType.KeyType, lvk))
	}

	var l9 string
	if dataType.ElemType.ByteSize == datatype.DynamicSize {
		l9 = fmt.Sprintf("%v_byte_size += %v;", varName, ig.ReadSizeExpression(*dataType.ElemType, lvv))
	}

	return generator.Lines(ls, l8, l9, "}")
}

func (g CxxMapGenerator) ToFuncMap() template.FuncMap {
	//TODO implement me
	panic("implement me")
}
