package datatype

// TraverseTypeTree traverses the type tree starting from the given DataType,
// calling visitor whenever a DataType is encountered.
func TraverseTypeTree(t *DataType, visitor func(t *DataType)) {
	visitor(t)
	if t.KeyType != nil {
		TraverseTypeTree(t.KeyType, visitor)
	}
	if t.ElemType != nil {
		TraverseTypeTree(t.ElemType, visitor)
	}
}
