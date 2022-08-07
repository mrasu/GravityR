package viewmodel

import (
	"strconv"
	"time"
)

type VmTimestamp time.Time

func (t VmTimestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).UnixMilli(), 10)), nil
}

type VmTimeDbLoad struct {
	Timestamp VmTimestamp `json:"timestamp"`
	Databases []*VmDbLoad `json:"databases"`
}

func NewVmTimeDbLoad(t time.Time) *VmTimeDbLoad {
	return &VmTimeDbLoad{
		Timestamp: VmTimestamp(t),
		Databases: []*VmDbLoad{},
	}
}
