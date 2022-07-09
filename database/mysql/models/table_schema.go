package models

type TableSchema struct {
	Name        string
	Columns     []*ColumnSchema
	PrimaryKeys []string
}
