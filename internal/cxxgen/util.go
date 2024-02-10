package cxxgen

import "nanoc/internal/datatype"

func isTriviallyCopiable(dataType datatype.DataType) bool {
	switch dataType.Kind {
	case datatype.Int8, datatype.Int32, datatype.Int64, datatype.Bool, datatype.Double:
		return true
	default:
		return false
	}
}
