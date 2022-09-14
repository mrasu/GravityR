package db_models

import "fmt"

type ColumnSchema struct {
	Name string
}

func (cs *ColumnSchema) String() string {
	return fmt.Sprintf(
		"ColumnSchema(name: %s)",
		cs.Name,
	)
}
