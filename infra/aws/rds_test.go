package aws_test

import (
	aAws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jarcoal/httpmock"
	"github.com/mrasu/GravityR/infra/aws"
	"github.com/mrasu/GravityR/infra/aws/models"
	"github.com/mrasu/GravityR/thelper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestRds_GetDBs(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	expDBs := []*models.RdsDB{
		{
			InstanceIdentifier: "gravityr1",
			DbiResourceId:      "db-XXXXXXXXXXXXXXXXXXXXXXXXX1",
		},
		{
			InstanceIdentifier: "gravityr2",
			DbiResourceId:      "db-XXXXXXXXXXXXXXXXXXXXXXXXX2",
		},
	}
	expReqBody := map[string][]string{
		"Action":                          {"DescribeDBInstances"},
		"Filters.Filter.1.Name":           {"engine"},
		"Filters.Filter.1.Values.Value.1": {"mysql"},
		"Version":                         {"2014-10-31"},
	}

	reqBody := map[string][]string{}
	httpmock.RegisterResponder("POST", "https://rds.ap-northeast-1.amazonaws.com/",
		func(req *http.Request) (*http.Response, error) {
			reqBody = unmarshalRdsBody(t, req)

			resp := thelper.ReadFromFiles(t, "aws/rds/DescribeDbInstances.xml")
			return httpmock.NewStringResponse(200, resp), nil
		},
	)

	rds := buildRds(t)
	dbs, err := rds.GetDBs([]string{"mysql"})
	assert.NoError(t, err)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	assert.Equal(t, expReqBody, reqBody)
	assert.Equal(t, expDBs, dbs)
}

func buildRds(t *testing.T) *aws.Rds {
	t.Helper()

	cfg := aAws.Config{
		Region: "ap-northeast-1",
	}
	cfg.HTTPClient = http.DefaultClient

	return aws.NewRds(cfg)
}

func unmarshalRdsBody(t *testing.T, req *http.Request) url.Values {
	t.Helper()

	err := req.ParseForm()
	assert.NoError(t, err)

	return req.Form
}
