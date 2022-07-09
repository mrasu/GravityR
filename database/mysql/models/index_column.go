package models

import "github.com/mrasu/GravityR/html/viewmodel"

type IndexColumn struct {
	Name string
}

func (ic *IndexColumn) ToViewModel() *viewmodel.VmIndexColumn {
	return &viewmodel.VmIndexColumn{Name: ic.Name}
}
