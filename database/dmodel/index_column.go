package dmodel

import (
	"github.com/mrasu/GravityR/html/viewmodel"
	"github.com/pkg/errors"
	"regexp"
)

type IndexColumn struct {
	Name string
}

func NewIndexColumn(name string) (*IndexColumn, error) {
	if !wordOnlyReg.MatchString(name) {
		return nil, errors.Errorf("Including non word character. %s", name)
	}

	return &IndexColumn{Name: name}, nil
}

func (ic *IndexColumn) ToViewModel() *viewmodel.VmIndexColumn {
	return &viewmodel.VmIndexColumn{Name: ic.Name}
}

var nonAsciiReg = regexp.MustCompile(`\W`)

func (ic *IndexColumn) SafeName() string {
	return nonAsciiReg.ReplaceAllString(ic.Name, "_")
}

func (ic *IndexColumn) Equals(other *IndexColumn) bool {
	return ic.Name == other.Name
}
