package main

import (
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/credentials"
	"github.com/awslabs/aws-sdk-go/service/sqs"
	"github.com/peter-edge/go-cacerts"
	"github.com/peter-edge/go-env"
	"github.com/peter-edge/go-tools/common"
)

type appEnv struct {
	AwsRegion string `env:"AWS_REGION,required"`
	QueueName string `env:"QUEUE_NAME,required"`
}

func main() {
	common.Main(do)
}

func do() error {
	var appEnv appEnv
	if err := env.Populate(&appEnv, env.PopulateOptions{}); err != nil {
		return err
	}
	s := sqs.New(
		&aws.Config{
			Credentials: credentials.NewEnvCredentials(),
			HTTPClient:  cacerts.NewHTTPClient(),
			Region:      appEnv.AwsRegion,
		},
	)
	createQueueOutput, err := s.CreateQueue(
		&sqs.CreateQueueInput{
			QueueName: aws.String(appEnv.QueueName),
		},
	)
	if err != nil {
		return err
	}
	_, err = s.PurgeQueue(
		&sqs.PurgeQueueInput{
			QueueURL: createQueueOutput.QueueURL,
		},
	)
	return err
}
