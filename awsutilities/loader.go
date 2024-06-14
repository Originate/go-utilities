package awsutilities

import (
	"context"
	"os"
	"reflect"

	"github.com/Originate/go-utilities/configutilities"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var awsVariables map[string]string = map[string]string{
	"AWS_ENDPOINT_URL":      "EndpointURL",
	"AWS_REGION":            "Region",
	"AWS_ACCESS_KEY_ID":     "AccessKeyID",
	"AWS_SECRET_ACCESS_KEY": "SecretAccessKey",
}

func LoadConfiguration(cfg configutilities.AWSSDKConfiguration) (aws.Config, error) {
	configStruct := reflect.ValueOf(cfg)
	for envKey, configKey := range awsVariables {
		if os.Getenv(envKey) == "" {
			os.Setenv(envKey, configStruct.FieldByName(configKey).String())
		}
	}

	return config.LoadDefaultConfig(context.Background())
}
