package cloud

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

var awsCfg *AWSConfig

type AWSConfig struct {
	Cfg aws.Config
}

// get aws config from profile
func GetAWSConfigWithProfile(profile string) (*AWSConfig, error) {
	if awsCfg != nil {
		return awsCfg, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}
	otelaws.AppendMiddlewares(&cfg.APIOptions)

	awsCfg = &AWSConfig{
		Cfg: cfg,
	}

	return awsCfg, nil
}

func GetAWSConfig() (*AWSConfig, error) {
	if awsCfg != nil {
		return awsCfg, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	awsCfg = &AWSConfig{cfg}

	return awsCfg, nil
}
