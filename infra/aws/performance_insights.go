package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	aPi "github.com/aws/aws-sdk-go-v2/service/pi"
	"github.com/aws/aws-sdk-go-v2/service/pi/types"
	"github.com/mrasu/GravityR/infra/aws/models"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"time"
)

const loadAvgMetricName = "db.load.avg"
const sqlGroupName = "db.sql"
const tokenizedSqlGroupName = "db.sql_tokenized"

type PerformanceInsights struct {
	client *aPi.Client
}

func NewPerformanceInsights(cfg aws.Config) *PerformanceInsights {
	return &PerformanceInsights{
		client: aPi.NewFromConfig(cfg),
	}
}

func (pi *PerformanceInsights) GetHalfDaySqlMetrics(db *models.RdsDB, start time.Time) ([]*models.PiSQLLoadAvg, error) {
	output, err := pi.getHalfDayMetrics(db, start, sqlGroupName)
	if err != nil {
		return nil, err
	}

	res := pi.convertResourceMetricsToAvgs(db, output, "db.sql.statement", "db.sql.tokenized_id")
	return res, nil
}

func (pi *PerformanceInsights) GetHalfDayTokenizedSqlMetrics(db *models.RdsDB, start time.Time) ([]*models.PiSQLLoadAvg, error) {
	output, err := pi.getHalfDayMetrics(db, start, tokenizedSqlGroupName)
	if err != nil {
		return nil, err
	}

	res := pi.convertResourceMetricsToAvgs(db, output, "db.sql_tokenized.statement", "db.sql_tokenized.id")
	return res, nil
}

func (pi *PerformanceInsights) getHalfDayMetrics(db *models.RdsDB, start time.Time, group string) (*aPi.GetResourceMetricsOutput, error) {
	output, err := pi.client.GetResourceMetrics(context.Background(), &aPi.GetResourceMetricsInput{
		Identifier: &db.DbiResourceId,
		MetricQueries: []types.MetricQuery{
			{
				Metric: lib.Ptr(loadAvgMetricName),
				GroupBy: &types.DimensionGroup{
					Group: &group,
					Limit: lib.Ptr(int32(10)),
				},
			},
		},
		ServiceType:     types.ServiceTypeRds,
		StartTime:       &start,
		EndTime:         lib.Ptr(start.Add(12 * time.Hour)),
		PeriodInSeconds: lib.Ptr(int32(300)),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetResourceMetrics of PerformanceInsights")
	}

	return output, nil
}

func (pi *PerformanceInsights) convertResourceMetricsToAvgs(db *models.RdsDB, output *aPi.GetResourceMetricsOutput, sqlKey, tokenizedIdKey string) []*models.PiSQLLoadAvg {
	var res []*models.PiSQLLoadAvg
	for _, m := range output.MetricList {
		if m.Key.Metric == nil || *m.Key.Metric != loadAvgMetricName {
			continue
		}
		if m.Key.Dimensions == nil {
			continue
		}

		sla := models.NewPiSQLLoadAvg(
			db.InstanceIdentifier,
			m.Key.Dimensions[sqlKey],
			m.Key.Dimensions[tokenizedIdKey],
		)
		res = append(res, sla)

		for _, dp := range m.DataPoints {
			if dp.Timestamp == nil {
				continue
			}

			if dp.Value == nil {
				continue
			}
			if *dp.Value == 0 {
				continue
			}
			sla.Values = append(sla.Values, &models.PiTimeValue{
				Value: *dp.Value,
				Time:  *dp.Timestamp,
			})
		}
	}

	return res
}
