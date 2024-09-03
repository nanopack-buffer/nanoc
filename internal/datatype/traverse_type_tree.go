package datatype

// TraverseTypeTree traverses the type tree starting from the given DataType,
// calling visitor whenever a DataType is encountered.
func TraverseTypeTree(t *DataType, visitor func(t *DataType) error) error {
	err := visitor(t)
	if err != nil {
		return err
	}
	if t.KeyType != nil {
		return TraverseTypeTree(t.KeyType, visitor)
	}
	if t.ElemType != nil {
		return TraverseTypeTree(t.ElemType, visitor)
	}
	return nil
}
