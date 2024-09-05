package cxxgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
	"strings"

	"github.com/iancoleman/strcase"
)

type mapGenerator struct {
	gm cxxCodeFragmentGeneratorMap
}

func (g mapGenerator) TypeDeclaration(dataType datatype.DataType) string {
	kg := g.gm[dataType.KeyType.Kind]
	ig := g.gm[dataType.ElemType.Kind]
	return fmt.Sprintf("std::unordered_map<%v, %v>", kg.TypeDeclaration(*dataType.KeyType), ig.TypeDeclaration(*dataType.ElemType))
}

func (g mapGenerator) ParameterDeclaration(dataType datatype.DataType, paramName string) string {
	return "const " + g.TypeDeclaration(dataType) + " &" + paramName
}

func (g mapGenerator) RValue(dataType datatype.DataType, argName string) string {
	return argName
}

func (g mapGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	if dataType.KeyType.ByteSize != datatype.DynamicSize && dataType.ElemType.ByteSize != datatype.DynamicSize {
		return fmt.Sprintf("%v.size() * %d", varName, dataType.KeyType.ByteSize+dataType.ElemType.ByteSize)
	}
	return varName + "_byte_size"
}

func (g mapGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name)
}

func (g mapGenerator) FieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	return fmt.Sprintf("%v(std::move(%v))", s, s)
}

func (g mapGenerator) FieldDeclaration(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name) + ";"
}

func (g mapGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)
	kg := g.gm[field.Type.KeyType.Kind]
	ig := g.gm[field.Type.ElemType.Kind]

	mapSizeVar := s + "_map_size"
	var readMapSize string
	if field.Type.KeyType.ByteSize != datatype.DynamicSize && field.Type.ElemType.ByteSize != datatype.DynamicSize {
		ctx.AddVariableToScope(mapSizeVar)
		// for maps with fixed size entries, the number of entries in the map can be calculated.
		readMapSize = fmt.Sprintf("const uint32_t %v = %v_byte_size / %d;", mapSizeVar, s, field.Type.KeyType.ByteSize+field.Type.ElemType.ByteSize)
	} else {
		readMapSize = generator.Lines(
			fmt.Sprintf("uint32_t %v;", mapSizeVar),
			fmt.Sprintf("reader.read_uint32(ptr, %v);", mapSizeVar),
			"ptr += 4;")
	}

	lv := ctx.NewLoopVar()
	ctx.AddVariableToScope(lv + "_value")

	ls := generator.Lines(
		readMapSize,
		fmt.Sprintf("%v.reserve(%v);", s, mapSizeVar),
		fmt.Sprintf("for (int %v = 0; %v < %v; %v++) {", lv, lv, mapSizeVar, lv),
		// read map key value from buffer
		kg.ReadValueFromBuffer(*field.Type.KeyType, lv+"_key", ctx),
		// create map entry
		fmt.Sprintf("    auto &%v_value = %v[%v_key];", lv, s, lv),
		// read value to entry
		ig.ReadValueFromBuffer(*field.Type.ElemType, lv+"_value", ctx),
		"}",
	)

	ctx.RemoveVariableFromScope(lv)
	ctx.RemoveVariableFromScope(lv + "_key")
	ctx.RemoveVariableFromScope(lv + "_value")

	return ls
}

func (g mapGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	i32g := g.gm[datatype.Int32]
	kg := g.gm[dataType.KeyType.Kind]
	ig := g.gm[dataType.ElemType.Kind]
	mapSizeVar := varName + "_map_size"

	// If the number of elements in the map is not read previously,
	// generate code to read it here.
	var l0 string
	if !ctx.IsVariableInScope(mapSizeVar) {
		l0 = i32g.ReadValueFromBuffer(*datatype.FromKind(datatype.Int32), mapSizeVar, ctx)
	}

	var l1 string
	if ctx.IsVariableInScope(varName) {
		l1 = fmt.Sprintf("%v = %v()", varName, g.TypeDeclaration(dataType))
	} else {
		l1 = fmt.Sprintf("%v %v;", varName, g.TypeDeclaration(dataType))
	}

	lv := ctx.NewLoopVar()
	ctx.AddVariableToScope(lv)
	ctx.AddVariableToScope(lv + "_value")

	ls := generator.Lines(
		l0,
		l1,
		fmt.Sprintf("%v.reserve(%v);", varName, mapSizeVar),
		fmt.Sprintf("for (int %v = 0; %v < %v; %v++) {", lv, lv, mapSizeVar, lv),
		// read map key value from buffer
		kg.ReadValueFromBuffer(*dataType.KeyType, lv+"_key", ctx),
		// create map entry
		fmt.Sprintf("    auto &%v_value = %v[%v_key];", lv, varName, lv),
		// read value to entry
		ig.ReadValueFromBuffer(*dataType.ElemType, lv+"_value", ctx),
		"}")

	ctx.RemoveVariableFromScope(lv)

	return ls
}

func (g mapGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
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
			fmt.Sprintf("const int32_t %v_byte_size = %v.size() * %d;", s, s, field.Type.ElemType.ByteSize+field.Type.KeyType.ByteSize),
			fmt.Sprintf("writer.write_field_size(%d, %v_byte_size, offset, buf);", field.Number, s),
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
		fmt.Sprintf("writer.write_field_size(%d, %v_byte_size, offset);", field.Number, s))
}

func (g mapGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
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
		fmt.Sprintf("auto %v = %v.first;", lvk, lv),
		fmt.Sprintf("auto %v = %v.second;", lvv, lv),
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

	ctx.RemoveVariableFromScope(lv)
	ctx.RemoveVariableFromScope(lvk)
	ctx.RemoveVariableFromScope(lvv)

	return generator.Lines(ls, l8, l9, "}")
}
