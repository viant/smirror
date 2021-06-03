package endpoint

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/viant/afs/option"
	"os"
)

const (
	awsRegionEnvKey = "AWS_REGION"
	awsCredentials  = "AWS_CREDENTIALS"
)

func GetAwsConfig(credLocation string) (*aws.Config, error) {
	var config *aws.Config
	var err error
	if credLocation == "" {
		credLocation = os.Getenv(awsCredentials)
	}
	if credLocation != "" {
		authConfig, err := NewAuthConfig(&option.Location{Path: credLocation})
		if err == nil {
			config, err = authConfig.AwsConfig()
		}
	} else {
		config = &aws.Config{}
	}
	if config != nil {
		if awsRegion := os.Getenv(awsRegionEnvKey); awsRegion != "" {
			config.Region = &awsRegion
		}
	}
	return config, err
}
