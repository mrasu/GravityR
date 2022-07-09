package viewmodel

type VmIndexTarget struct {
	TableName string           `json:"tableName"`
	Columns   []*VmIndexColumn `json:"columns"`
}
