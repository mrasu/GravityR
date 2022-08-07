package models

import (
	"github.com/mrasu/GravityR/html/viewmodel"
	"github.com/mrasu/GravityR/lib"
	"time"
)

type PiSQLLoadAvg struct {
	DbName      string
	SQL         string
	TokenizedId string
	Values      []*PiTimeValue
}

func NewPiSQLLoadAvg(dbName, sql, tokenizedId string) *PiSQLLoadAvg {
	return &PiSQLLoadAvg{
		DbName:      dbName,
		SQL:         sql,
		TokenizedId: tokenizedId,
	}
}

func ConvertPiSQLLoadAvgsToVms(start, end time.Time, avgs []*PiSQLLoadAvg) []*viewmodel.VmTimeDbLoad {
	sqlLoadMap := map[time.Time]*viewmodel.VmTimeDbLoad{}
	dbLoadMap := map[time.Time]map[string]*viewmodel.VmDbLoad{}
	for _, avg := range avgs {
		for _, v := range avg.Values {
			baseTime := lib.NormalizeTimeByHour(v.Time)
			if _, ok := sqlLoadMap[baseTime]; !ok {
				sqlLoadMap[baseTime] = viewmodel.NewVmTimeDbLoad(baseTime)
			}
			if _, ok := dbLoadMap[baseTime]; !ok {
				dbLoadMap[baseTime] = map[string]*viewmodel.VmDbLoad{}
			}
			load := sqlLoadMap[baseTime]

			if _, ok := dbLoadMap[baseTime][avg.DbName]; !ok {
				dbLoad := viewmodel.NewVmDbLoad(avg.DbName)
				dbLoadMap[baseTime][avg.DbName] = dbLoad
				load.Databases = append(load.Databases, dbLoad)
			}
			dbLoad := dbLoadMap[baseTime][avg.DbName]

			found := false
			for _, loadOfSql := range dbLoad.Sqls {
				if loadOfSql.Sql == avg.SQL {
					if loadOfSql.LoadMax < v.Value {
						loadOfSql.LoadMax = v.Value
					}
					loadOfSql.LoadSum += v.Value

					found = true
					break
				}
			}

			if !found {
				dbLoad.Sqls = append(dbLoad.Sqls, &viewmodel.VmDbLoadOfSql{
					Sql:         avg.SQL,
					LoadMax:     v.Value,
					LoadSum:     v.Value,
					TokenizedId: avg.TokenizedId,
				})
			}
		}
	}

	target := lib.NormalizeTimeByHour(start).UTC()
	var res []*viewmodel.VmTimeDbLoad
	for !target.After(end) {
		if l, ok := sqlLoadMap[target]; ok {
			res = append(res, l)
		} else {
			res = append(res, viewmodel.NewVmTimeDbLoad(target))
		}
		target = target.Add(1 * time.Hour)
	}

	return res
}
