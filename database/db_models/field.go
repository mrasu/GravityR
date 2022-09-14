package db_models

type Field struct {
	AsName  string
	Columns []*FieldColumn
}

func (f *Field) Name() string {
	if f.AsName != "" {
		return f.AsName
	} else if len(f.Columns) == 1 {
		return f.Columns[0].Name
	}

	return ""
}
