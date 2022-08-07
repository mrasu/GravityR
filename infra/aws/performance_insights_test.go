package aws_test

import (
	"encoding/json"
	aAws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jarcoal/httpmock"
	"github.com/mrasu/GravityR/infra/aws"
	"github.com/mrasu/GravityR/infra/aws/models"
	"github.com/mrasu/GravityR/thelper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type awsPiReqBody struct {
	PeriodInSeconds int64
	StartTime       int64
	EndTime         int64
	MetricQueries   []awsPiMetricQuery
}

type awsPiMetricQuery struct {
	Metric  string
	GroupBy awsPiGroupBy
}

type awsPiGroupBy struct {
	Group string
	Limit int64
}

func TestPerformanceInsights_GetHalfDaySqlMetrics(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	expReqBody := awsPiReqBody{
		PeriodInSeconds: 300,
		StartTime:       1659072300,
		EndTime:         1659115500,
		MetricQueries: []awsPiMetricQuery{
			{
				Metric: "db.load.avg",
				GroupBy: awsPiGroupBy{
					Group: "db.sql",
					Limit: 10,
				},
			},
		},
	}

	expAvgs := []*models.PiSQLLoadAvg{
		{
			DbName:      "db-x",
			SQL:         "INSERT INTO todos(user_id...",
			TokenizedId: "70FB8D1F1B692FE0C458211B3EE75B5710E81F1C",
			Values: []*models.PiTimeValue{
				{
					Time:  time.Unix(1659073200, 0).In(time.UTC),
					Value: 0.18661971830985916,
				},
			},
		},
		{
			DbName:      "db-x",
			SQL:         "explain analyze SELECT name...",
			TokenizedId: "64E6FE2C360536B8BED607E6A6B66BC463667A4A",
			Values: []*models.PiTimeValue{
				{
					Time:  time.Unix(1659073500, 0).In(time.UTC),
					Value: 0.078125,
				},
			},
		},
		{
			DbName:      "db-x",
			SQL:         "SELECT name, t.description...",
			TokenizedId: "D51E26E104C8374C3547D9AEC929104FB6A2D9AE",
			Values: []*models.PiTimeValue{
				{
					Time:  time.Unix(1659089400, 0).In(time.UTC),
					Value: 0.012355,
				},
			},
		},
	}

	reqBody := awsPiReqBody{}
	httpmock.RegisterResponder("POST", "https://pi.ap-northeast-1.amazonaws.com/",
		func(req *http.Request) (*http.Response, error) {
			reqBody = unmarshalPIBody(t, req)

			sqlContent := thelper.ReadFromFiles(t, "aws/pi/GetResourceMetrics_sql.json")
			return httpmock.NewStringResponse(200, sqlContent), nil
		},
	)

	pi := buildPI(t)
	avgs, err := pi.GetHalfDaySqlMetrics(&models.RdsDB{
		InstanceIdentifier: "db-x",
		DbiResourceId:      "dbi",
	}, time.Unix(1659072300, 0))

	assert.NoError(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	assert.Equal(t, expReqBody, reqBody)
	assert.Equal(t, expAvgs, avgs)
}

func TestPerformanceInsights_GetHalfDayTokenizedSqlMetrics(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	expReqBody := awsPiReqBody{
		PeriodInSeconds: 300,
		StartTime:       1659072300,
		EndTime:         1659115500,
		MetricQueries: []awsPiMetricQuery{
			{
				Metric: "db.load.avg",
				GroupBy: awsPiGroupBy{
					Group: "db.sql_tokenized",
					Limit: 10,
				},
			},
		},
	}

	expAvgs := []*models.PiSQLLoadAvg{
		{
			DbName:      "db-x",
			SQL:         "INSERT INTO `todos`(`user_id`...",
			TokenizedId: "70FB8D1F1B692FE0C458211B3EE75B5710E81F1C",
			Values: []*models.PiTimeValue{
				{
					Time:  time.Unix(1659073200, 0).In(time.UTC),
					Value: 0.18661971830985916,
				},
			},
		},
		{
			DbName:      "db-x",
			SQL:         "explain analyze SELECT NAME...",
			TokenizedId: "64E6FE2C360536B8BED607E6A6B66BC463667A4A",
			Values: []*models.PiTimeValue{
				{
					Time:  time.Unix(1659073500, 0).In(time.UTC),
					Value: 0.078125,
				},
			},
		},
		{
			DbName:      "db-x",
			SQL:         "SELECT NAME, t.`description`...",
			TokenizedId: "D51E26E104C8374C3547D9AEC929104FB6A2D9AE",
			Values: []*models.PiTimeValue{
				{
					Time:  time.Unix(1659089400, 0).In(time.UTC),
					Value: 0.012355,
				},
			},
		},
	}

	reqBody := awsPiReqBody{}
	httpmock.RegisterResponder("POST", "https://pi.ap-northeast-1.amazonaws.com/",
		func(req *http.Request) (*http.Response, error) {
			reqBody = unmarshalPIBody(t, req)

			sqlContent := thelper.ReadFromFiles(t, "aws/pi/GetResourceMetrics_tokenized_sql.json")
			return httpmock.NewStringResponse(200, sqlContent), nil
		},
	)

	pi := buildPI(t)
	avgs, err := pi.GetHalfDayTokenizedSqlMetrics(&models.RdsDB{
		InstanceIdentifier: "db-x",
		DbiResourceId:      "dbi",
	}, time.Unix(1659072300, 0))

	assert.NoError(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	assert.Equal(t, expReqBody, reqBody)
	assert.Equal(t, expAvgs, avgs)
}

func buildPI(t *testing.T) *aws.PerformanceInsights {
	t.Helper()

	cfg := aAws.Config{
		Region: "ap-northeast-1",
	}

	cfg.HTTPClient = http.DefaultClient

	return aws.NewPerformanceInsights(cfg)
}

func unmarshalPIBody(t *testing.T, req *http.Request) awsPiReqBody {
	t.Helper()

	out, err := ioutil.ReadAll(req.Body)
	assert.NoError(t, err)
	reqBody := awsPiReqBody{}
	err = json.Unmarshal(out, &reqBody)
	assert.NoError(t, err)

	return reqBody
}
