package dmodel

import (
	"fmt"
	"github.com/mrasu/GravityR/html/viewmodel"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"regexp"
	"strings"
)

type IndexTarget struct {
	TableName string
	Columns   []*IndexColumn
}

var indexTargetReg = regexp.MustCompile(`(.+?):((.+?(\+)?)+)`)
var wordOnlyReg = regexp.MustCompile(`\A\w+\z`)

const indexTargetColumnSeparator = "+"

func NewIndexTarget(tableName string, columns []string) *IndexTarget {
	return &IndexTarget{
		TableName: tableName,
		Columns:   lo.Map(columns, func(col string, _ int) *IndexColumn { return &IndexColumn{Name: col} }),
	}
}

func NewIndexTargetFromText(text string) (*IndexTarget, error) {
	m := indexTargetReg.FindStringSubmatch(text)
	if m == nil {
		return nil, errors.Errorf("Not appropriate text for index: %s", text)
	}

	tName := m[1]
	if !wordOnlyReg.MatchString(tName) {
		return nil, errors.Errorf("Including non word character. %s", tName)
	}

	it := &IndexTarget{TableName: tName}
	columns := strings.Split(m[2], indexTargetColumnSeparator)
	for _, c := range columns {
		ic, err := NewIndexColumn(c)
		if err != nil {
			return nil, err
		}

		it.Columns = append(it.Columns, ic)
	}
	return it, nil
}

func (it *IndexTarget) ToViewModel() *viewmodel.VmIndexTarget {
	vm := &viewmodel.VmIndexTarget{
		TableName: it.TableName,
	}
	for _, c := range it.Columns {
		vm.Columns = append(vm.Columns, c.ToViewModel())
	}

	return vm
}

func (it *IndexTarget) CombinationString() string {
	txt := fmt.Sprintf(
		"table: %s, columns: [%s])",
		it.TableName,
		lib.Join(it.Columns, ", ", func(f *IndexColumn) string { return f.Name }),
	)
	return txt
}

func (it *IndexTarget) Equals(other *IndexTarget) bool {
	if it.TableName != other.TableName {
		return false
	}
	if len(it.Columns) != len(other.Columns) {
		return false
	}

	for i, col := range it.Columns {
		if !col.Equals(other.Columns[i]) {
			return false
		}
	}

	return true
}

func (it *IndexTarget) IsSafe() bool {
	if !wordOnlyReg.MatchString(it.TableName) {
		return false
	}

	for _, c := range it.Columns {
		if c.Name != c.SafeName() {
			return false
		}
	}

	return true
}

func (it *IndexTarget) HasSameIdxColumns(info *IndexTarget) bool {
	if len(info.Columns) != len(it.Columns) {
		return false
	}

	for i, col := range info.Columns {
		if !strings.EqualFold(it.Columns[i].Name, col.Name) {
			return false
		}
	}

	return true
}
