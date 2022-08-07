package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
)

func NewAwsConfig(endpointURL string, verbose bool) (aws.Config, error) {
	var optFns []func(*config.LoadOptions) error

	if endpointURL != "" {
		optFns = append(optFns, config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL: endpointURL,
					}, nil
				},
			),
		))
	}
	if verbose {
		optFns = append(optFns, config.WithClientLogMode(aws.LogRetries|aws.LogRequestWithBody|aws.LogResponseWithBody))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), optFns...)
	if err != nil {
		return cfg, errors.Wrap(err, "failed to load aws config")
	}

	return cfg, nil
}
