package npschema

type Schema interface {
	isSchema()
}

type PartialSchema interface {
	isPartialSchema()
}
