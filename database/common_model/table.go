package common_model

type Table struct {
	AsName    string
	Name      string
	IsLateral bool
}

func (t *Table) AsOrName() string {
	if t.AsName != "" {
		return t.AsName
	} else {
		return t.Name
	}
}
