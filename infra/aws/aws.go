package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/smithy-go/logging"
	"github.com/pkg/errors"
	"os"
)

const mockEndpoint = "http://localhost:8080"
const verboseMode = aws.LogRetries | aws.LogRequest | aws.LogResponse

func NewAwsConfig(useMock bool, verbose bool) (aws.Config, error) {
	if useMock {
		return newMockedConfig(verbose), nil
	}

	var optFns []func(*config.LoadOptions) error
	if verbose {
		optFns = append(optFns, config.WithClientLogMode(verboseMode))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), optFns...)
	if err != nil {
		return cfg, errors.Wrap(err, "failed to load aws config")
	}

	return cfg, nil
}

func newMockedConfig(verbose bool) aws.Config {
	cfg := aws.Config{
		Logger: logging.NewStandardLogger(os.Stderr),
	}

	cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: mockEndpoint,
			}, nil
		},
	)
	if verbose {
		cfg.ClientLogMode = verboseMode
	}

	return cfg
}
