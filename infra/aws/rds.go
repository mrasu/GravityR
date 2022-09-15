package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	aRds "github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
)

type Rds struct {
	client *aRds.Client
}

func NewRds(cfg aws.Config) *Rds {
	cli := aRds.NewFromConfig(cfg)
	return &Rds{
		client: cli,
	}
}

func (rds *Rds) GetDBs(engines []string) ([]*RdsDB, error) {
	output, err := rds.client.DescribeDBInstances(context.Background(), &aRds.DescribeDBInstancesInput{
		Filters: []types.Filter{
			{Name: lib.Ptr("engine"), Values: engines},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to describe db instances")
	}

	var res []*RdsDB
	for _, ins := range output.DBInstances {
		res = append(res, &RdsDB{
			InstanceIdentifier: *ins.DBInstanceIdentifier,
			DbiResourceId:      *ins.DbiResourceId,
		})
	}

	return res, nil
}
