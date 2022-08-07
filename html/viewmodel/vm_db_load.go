package viewmodel

type VmDbLoad struct {
	Name string           `json:"name"`
	Sqls []*VmDbLoadOfSql `json:"sqls"`
}

func NewVmDbLoad(name string) *VmDbLoad {
	return &VmDbLoad{
		Name: name,
		Sqls: []*VmDbLoadOfSql{},
	}
}
