package commonaws

import (
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/credentials"
	"github.com/peter-edge/go-cacerts"
	"github.com/peter-edge/go-env"
	"github.com/peter-edge/go-tools/common"
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
