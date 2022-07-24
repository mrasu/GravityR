package models

import (
	"fmt"
	"github.com/mrasu/GravityR/html/viewmodel"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type IndexTarget struct {
	TableName string
	Columns   []*IndexColumn
}

var indexTargetReg = regexp.MustCompile("(.+?):((.+?(\\\\+)?)+)")
var indexTargetWithBacktickReg = regexp.MustCompile("(.+?)`:((.+?(`\\\\+)?)+)")

const indexTargetColumnSeparator = "+"
const indexTargetColumnWithBacktickSeparator = "`+"

func ParseIndexTarget(text string) (*IndexTarget, error) {
	m := indexTargetReg.FindStringSubmatch(text)
	if m == nil {
		return nil, errors.Errorf("Not appropriate text for index: %s", text)
	}

	it := &IndexTarget{TableName: m[1]}
	columns := strings.Split(m[2], indexTargetColumnSeparator)
	for _, c := range columns {
		it.Columns = append(it.Columns, &IndexColumn{Name: c})
	}
	return it, nil
}

func ParseIndexTargetWithBacktick(text string) (*IndexTarget, error) {
	m := indexTargetWithBacktickReg.FindStringSubmatch(text)
	if m == nil {
		return nil, errors.Errorf("Not appropriate text for index: %s", text)
	}

	it := &IndexTarget{TableName: m[1]}
	columns := strings.Split(m[2], indexTargetColumnWithBacktickSeparator)
	for _, c := range columns {
		it.Columns = append(it.Columns, &IndexColumn{Name: c})
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

func (it *IndexTarget) String() string {
	txt := fmt.Sprintf(
		"IndexTargetTable(table: %s, columns: [%s])",
		it.TableName,
		lib.JoinF(it.Columns, ", ", func(f *IndexColumn) string { return f.Name }),
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
