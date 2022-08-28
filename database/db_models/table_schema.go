package db_models

type TableSchema struct {
	Name        string
	Columns     []*ColumnSchema
	PrimaryKeys []string
}

func CreateTableSchemas[T any](tables []string, vals []T, f func(T) (string, string, bool)) []*TableSchema {
	ts := map[string]*TableSchema{}
	for _, v := range vals {
		var t *TableSchema
		tName, colName, isPK := f(v)
		if et, ok := ts[tName]; ok {
			t = et
		} else {
			t = &TableSchema{
				Name: tName,
			}
			ts[tName] = t
		}

		t.Columns = append(t.Columns, &ColumnSchema{
			Name: colName,
		})

		if isPK {
			t.PrimaryKeys = append(t.PrimaryKeys, colName)
		}
	}

	var res []*TableSchema
	for _, table := range tables {
		if t, ok := ts[table]; ok {
			res = append(res, t)
		} else {
			res = append(res, nil)
		}
	}

	return res
}
