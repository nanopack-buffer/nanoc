package cxxgen

import (
	"fmt"
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"

	"github.com/iancoleman/strcase"
)

type arrayGenerator struct {
	gm cxxCodeFragmentGeneratorMap
}

func (g arrayGenerator) TypeDeclaration(dataType datatype.DataType) string {
	ig := g.gm[dataType.ElemType.Kind]
	return fmt.Sprintf("std::vector<%v>", ig.TypeDeclaration(*dataType.ElemType))
}

func (g arrayGenerator) ParameterDeclaration(dataType datatype.DataType, paramName string) string {
	return fmt.Sprintf("const %v &%v", g.TypeDeclaration(dataType), paramName)
}

func (g arrayGenerator) RValue(dataType datatype.DataType, argName string) string {
	return argName
}

func (g arrayGenerator) ReadSizeExpression(dataType datatype.DataType, varName string) string {
	if dataType.ElemType.ByteSize != datatype.DynamicSize {
		return fmt.Sprintf("%v.size() * %d", varName, dataType.ElemType.ByteSize)
	}
	return fmt.Sprintf("%v_byte_size", varName)
}

func (g arrayGenerator) ConstructorFieldParameter(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name)
}

func (g arrayGenerator) FieldInitializer(field npschema.MessageField) string {
	s := strcase.ToSnake(field.Name)
	return fmt.Sprintf("%v(std::move(%v))", s, s)
}

func (g arrayGenerator) FieldDeclaration(field npschema.MessageField) string {
	return g.TypeDeclaration(field.Type) + " " + strcase.ToSnake(field.Name) + ";"
}

func (g arrayGenerator) ReadFieldFromBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)

	var itemType datatype.DataType
	if field.Type.Kind == datatype.Optional {
		itemType = *field.Type.ElemType.ElemType
	} else {
		itemType = *field.Type.ElemType
	}

	// array generator may receive optional array field
	// here we are unwrapping the optional type to obtain the underlying array type
	var fieldType datatype.DataType
	if field.Type.Kind == datatype.Optional {
		fieldType = *field.Type.ElemType
	} else {
		fieldType = field.Type
	}

	var vecInit string
	if field.Type.Kind == datatype.Optional {
		vecInit = fmt.Sprintf("%v = std::make_unique<%v>();", s, g.TypeDeclaration(*field.Type.ElemType))
	}

	var vecSizeDeclr string
	if itemType.ByteSize != datatype.DynamicSize {
		// for arrays with fixed size items, the number of elements in the array can be calculated.
		vecSizeDeclr = generator.Lines(
			fmt.Sprintf("const int32_t %v_byte_size = reader.read_field_size(%d);", s, field.Number),
			fmt.Sprintf("const int32_t %v_vec_size = %v_byte_size / %d;", s, s, itemType.ByteSize))
		ctx.AddVariableToScope(s + "_vec_size")
	}
	ctx.AddVariableToScope(s)
	return generator.Lines(
		vecInit,
		vecSizeDeclr,
		g.ReadValueFromBuffer(fieldType, s, ctx))
}

func (g arrayGenerator) ReadValueFromBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	i32g := g.gm[datatype.Int32]
	ig := g.gm[dataType.ElemType.Kind]
	vecSizeVar := varName + "_vec_size"

	lv := ctx.NewLoopVar()

	// If the number of elements in the vector is not read previously,
	// generate code to read it here.
	var readVecSize string
	if !ctx.IsVariableInScope(vecSizeVar) {
		readVecSize = i32g.ReadValueFromBuffer(*datatype.FromKind(datatype.Int32), vecSizeVar, ctx)
	}

	var vecResizeOrDeclr string
	if ctx.IsVariableInScope(varName) {
		vecResizeOrDeclr = fmt.Sprintf("%v.resize(%v);", varName, vecSizeVar)
	} else {
		vecResizeOrDeclr = fmt.Sprintf("%v %v(%v);", g.TypeDeclaration(dataType), varName, vecSizeVar)
	}

	readTo := lv + "_item"
	ctx.AddVariableToScope(readTo)

	ls := generator.Lines(
		readVecSize,
		vecResizeOrDeclr,
		fmt.Sprintf("for (int %v = 0; %v < %v; ++%v) {", lv, lv, vecSizeVar, lv),
		fmt.Sprintf("    auto &%v = %v[%v];", readTo, varName, lv),
		ig.ReadValueFromBuffer(*dataType.ElemType, readTo, ctx),
		"}")

	ctx.RemoveVariableFromScope(lv)
	ctx.RemoveVariableFromScope(readTo)

	return ls
}

func (g arrayGenerator) WriteFieldToBuffer(field npschema.MessageField, ctx generator.CodeContext) string {
	s := strcase.ToSnake(field.Name)

	if field.Type.ElemType.ByteSize != datatype.DynamicSize {
		// the array has fixed size elements, so the total size of the array
		// data can be calculated directly:
		//
		//     sizeof(element type) * number of elements in the vector +
		//
		// the number of elements in the array is not written to the buffer,
		// unlike elements with dynamic size, since it can be determined easily:
		//
		//     number of elements in the vector = total byte size of vector /
		//     sizeof(element type).
		lv := ctx.NewLoopVar()
		ig := g.gm[field.Type.ElemType.Kind]
		ls := generator.Lines(
			fmt.Sprintf("const int32_t %v_byte_size = %v.size() * %d;", s, s, field.Type.ElemType.ByteSize),
			fmt.Sprintf("writer.write_field_size(%d, %v_byte_size, offset);", field.Number, s),
			fmt.Sprintf("for (const auto &%v : %v) {", lv, s),
			ig.WriteVariableToBuffer(*field.Type.ElemType, lv, ctx),
			"}")
		ctx.RemoveVariableFromScope(lv)
		return ls
	}

	return generator.Lines(
		g.WriteVariableToBuffer(field.Type, s, ctx),
		fmt.Sprintf("writer.write_field_size(%d, %v_byte_size, offset);", field.Number, s))
}

func (g arrayGenerator) WriteVariableToBuffer(dataType datatype.DataType, varName string, ctx generator.CodeContext) string {
	ig := g.gm[dataType.ElemType.Kind]
	i32g := g.gm[datatype.Int32]
	vecSizeVar := varName + "_vec_size"
	isItemDynamicSize := dataType.ElemType.ByteSize == datatype.DynamicSize

	vecSizeDeclr := fmt.Sprintf("const size_t %v = %v.size();", vecSizeVar, varName)
	writeVecSizeStmt := i32g.WriteVariableToBuffer(*datatype.FromKind(datatype.Int32), vecSizeVar, ctx)

	var byteSizeDeclr string
	if isItemDynamicSize {
		// elements in the array are dynamically sized,
		// so the total byte size of the array cannot be determined statically.
		// here we declare a variable for storing the total byte size of all the
		// elements in the array, which is accumulated later in the loop
		byteSizeDeclr = fmt.Sprintf("int32_t %v_byte_size = sizeof(int32_t);", varName)
	}

	lv := ctx.NewLoopVar()
	ls := generator.Lines(
		vecSizeDeclr,
		writeVecSizeStmt,
		byteSizeDeclr,
		fmt.Sprintf("for (auto &%v : %v) {", lv, varName),
		ig.WriteVariableToBuffer(*dataType.ElemType, lv, ctx))

	var l5 string
	if isItemDynamicSize {
		l5 = fmt.Sprintf("%v_byte_size += %v;", varName, ig.ReadSizeExpression(*dataType.ElemType, lv))
	}

	return generator.Lines(
		ls,
		l5,
		"}")
}
