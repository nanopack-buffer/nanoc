package datatype

type Schema interface {
	SchemaPathAbsolute() string

	// DataType returns a DataType describing the data type this Schema defines.
	// Returns nil if this schema does not define a data type, for example a service schema.
	DataType() *DataType

	// SchemaName returns the name of the thing this schema defines.
	//   - For a message schema, this is the message name.
	//   - For an enum schema, this is the enum name.
	//   - For a service schema, this is the service name.
	SchemaName() string
}

type PartialSchema interface {
	IsPartialSchema()
}
