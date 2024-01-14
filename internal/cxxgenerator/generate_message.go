package cxxgenerator

import (
	"nanoc/internal/datatype"
	"nanoc/internal/generator"
	"nanoc/internal/npschema"
)

func GenerateMessageClass(msgSchema npschema.Message) error {
	ng := CxxNumberGenerator{}
	gm := generator.MessageCodeGeneratorMap{
		datatype.Int8:  ng,
		datatype.Int32: ng,
		datatype.Int64: ng,
	}
	gm[datatype.Array] = CxxArrayGenerator{gm}

	//TODO: implement me
	return nil
}
