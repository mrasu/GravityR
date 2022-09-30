package dig

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/pi"
	piTypes "github.com/aws/aws-sdk-go-v2/service/pi/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	rdsTypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/mrasu/GravityR/infra/aws"
	"github.com/mrasu/GravityR/lib"
	"github.com/mrasu/GravityR/thelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func Test_runPerformanceInsightsDig(t *testing.T) {
	start := time.Date(2001, 2, 3, 4, 0, 0, 0, time.UTC)
	pr := piRunner{
		start: start.Format(time.RFC3339),
	}

	m := piMock{}
	rdsCli := m.buildRdsClient()
	piCli := m.buildPiClient()

	thelper.InjectClientDist()
	thelper.CreateTemp(t, "tmp.html", func(tmpfile *os.File) {
		err := pr.dig(tmpfile.Name(), rdsCli, piCli)
		require.NoError(t, err)

		html, err := os.ReadFile(tmpfile.Name())
		require.NoError(t, err)

		assert.Contains(t, string(html), fmt.Sprintf(`sqlDbLoads":[{"timestamp":%d`, start.UnixMilli()))
		assert.Contains(t, string(html), `name":"dummy_instance","sqls":[{"sql":"SELECT * FROM users","loadMax":10.2`)
	})
}

type piMock struct{}

func (m *piMock) buildRdsClient() *aws.Rds {
	return aws.NewRdsWithCli(&rdsCliMock{})
}

func (m *piMock) buildPiClient() *aws.PerformanceInsights {
	return aws.NewPerformanceInsightsWithCli(&piCliMock{})
}

type rdsCliMock struct{}

func (r *rdsCliMock) DescribeDBInstances(context.Context, *rds.DescribeDBInstancesInput, ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error) {
	return &rds.DescribeDBInstancesOutput{
		DBInstances: []rdsTypes.DBInstance{
			{
				DBInstanceIdentifier: lib.Ptr("dummy_instance"),
				DbiResourceId:        lib.Ptr("dummy_resource"),
			},
		},
	}, nil
}

type piCliMock struct{}

func (r *piCliMock) GetResourceMetrics(context.Context, *pi.GetResourceMetricsInput, ...func(*pi.Options)) (*pi.GetResourceMetricsOutput, error) {
	return &pi.GetResourceMetricsOutput{
		MetricList: []piTypes.MetricKeyDataPoints{
			{
				Key: &piTypes.ResponseResourceMetricKey{
					Metric: lib.Ptr("db.load.avg"),
					Dimensions: map[string]string{
						"db.sql.statement":    "SELECT * FROM users",
						"db.sql.tokenized_id": "SELECT * FROM users",
					},
				},
				DataPoints: []piTypes.DataPoint{
					{
						Timestamp: lib.Ptr(time.Date(2001, 2, 3, 4, 0, 0, 0, time.UTC)),
						Value:     lib.Ptr(10.2),
					},
				},
			},
		},
	}, nil
}
