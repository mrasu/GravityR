package common_model

type FieldColumn struct {
	// internal name to link with parent scope
	ReferenceName string

	Table string
	Name  string
	Type  FieldType
}