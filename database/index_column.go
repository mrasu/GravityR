package database

import (
	"github.com/mrasu/GravityR/html/viewmodel"
	"github.com/pkg/errors"
	"regexp"
)

type IndexColumn struct {
	name string
}

func NewIndexColumn(name string) (*IndexColumn, error) {
	if !wordOnlyReg.MatchString(name) {
		return nil, errors.Errorf("Including non word character. %s", name)
	}

	return &IndexColumn{name: name}, nil
}

func (ic *IndexColumn) ToViewModel() *viewmodel.VmIndexColumn {
	return &viewmodel.VmIndexColumn{Name: ic.name}
}

var nonAsciiReg = regexp.MustCompile(`\W`)

func (ic *IndexColumn) SafeName() string {
	return nonAsciiReg.ReplaceAllString(ic.name, "_")
}

func (ic *IndexColumn) Equals(other *IndexColumn) bool {
	return ic.name == other.name
}
