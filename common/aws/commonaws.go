package commonaws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/codeship/go-tools/common"
	"github.com/peter-edge/go-cacerts"
	"github.com/peter-edge/go-env"
)

type awsEnv struct {
	AwsRegion string `env:"AWS_REGION,required"`
}

func Main(do func(*aws.Config) error) {
	common.Main(func() error {
		var awsEnv awsEnv
		if err := env.Populate(&awsEnv, env.PopulateOptions{}); err != nil {
			return err
		}
		return do(
			&aws.Config{
				Credentials: credentials.NewEnvCredentials(),
				HTTPClient:  cacerts.NewHTTPClient(),
				Region:      awsEnv.AwsRegion,
			},
		)
	})
}
