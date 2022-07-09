package viewmodel

import (
	"path/filepath"
	"strings"
)

type VmExaminationCommandOption struct {
	IsShort bool   `json:"isShort"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}

func CreateOutputExaminationOption(addPrefix bool, filename string) *VmExaminationCommandOption {
	val := filename
	if addPrefix {
		val = strings.TrimSuffix(filename, filepath.Ext(filename)) + "_examine" + filepath.Ext(filename)
	}

	return &VmExaminationCommandOption{
		IsShort: true,
		Name:    "o",
		Value:   val,
	}
}
