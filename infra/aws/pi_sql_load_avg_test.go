package aws_test

import (
	"github.com/mrasu/GravityR/html/viewmodel"
	"github.com/mrasu/GravityR/infra/aws"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConvertPiSQLLoadAvgsToVms(t *testing.T) {
	start := time.Unix(1659074400, 0).UTC()
	end := start.Add(59 * time.Minute)

	tests := []struct {
		name   string
		avgs   []*aws.PiSQLLoadAvg
		expVms []*viewmodel.VmTimeDbLoad
	}{
		{
			name: "query same time",
			avgs: []*aws.PiSQLLoadAvg{{
				DbName:      "db-x",
				SQL:         "SELECT * FROM tasks...",
				TokenizedId: "aaa",
				Values: []*aws.PiTimeValue{
					{
						Time:  start,
						Value: 0.1,
					},
					{
						Time:  start,
						Value: 0.4,
					},
					{
						Time:  start,
						Value: 0.2,
					},
				},
			}},
			expVms: []*viewmodel.VmTimeDbLoad{{
				Timestamp: viewmodel.VmTimestamp(start),
				Databases: []*viewmodel.VmDbLoad{{
					Name: "db-x",
					Sqls: []*viewmodel.VmDbLoadOfSql{{
						Sql:         "SELECT * FROM tasks...",
						LoadMax:     0.4,
						LoadSum:     0.7,
						TokenizedId: "aaa",
					}},
				}},
			}},
		},
		{
			name: "query different period",
			avgs: []*aws.PiSQLLoadAvg{{
				DbName:      "db-x",
				SQL:         "SELECT * FROM tasks...",
				TokenizedId: "aaa",
				Values: []*aws.PiTimeValue{
					{
						Time:  start,
						Value: 0.1,
					},
					{
						Time:  start.Add(10 * time.Minute),
						Value: 0.4,
					},
					{
						Time:  start.Add(20 * time.Minute),
						Value: 0.2,
					},
				},
			}},
			expVms: []*viewmodel.VmTimeDbLoad{{
				Timestamp: viewmodel.VmTimestamp(start),
				Databases: []*viewmodel.VmDbLoad{{
					Name: "db-x",
					Sqls: []*viewmodel.VmDbLoadOfSql{{
						Sql:         "SELECT * FROM tasks...",
						LoadMax:     0.4,
						LoadSum:     0.7,
						TokenizedId: "aaa",
					}},
				}},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualVms := aws.ConvertPiSQLLoadAvgsToVms(start, end, tt.avgs)
			assert.Equal(t, tt.expVms, actualVms)
		})
	}
}
