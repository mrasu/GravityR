package db_models

type FieldType int

const (
	FieldReference FieldType = iota
	FieldCondition
	FieldAggregation
	FieldSubquery
	FieldStar
)
