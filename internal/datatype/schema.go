package datatype

type Schema interface {
	SchemaPathAbsolute() string

	DataType() DataType
}

type PartialSchema interface {
	IsPartialSchema()
}
