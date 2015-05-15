package main

import (
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/credentials"
	"github.com/awslabs/aws-sdk-go/service/sqs"
	"github.com/peter-edge/go-appenv"
	"github.com/peter-edge/go-cacerts"
	"github.com/peter-edge/go-tools/common"
)

var (
	allEnvKeys = []string{
		"AWS_REGION",
		"QUEUE_NAME",
	}
)

type env struct {
	AwsRegion string `env:"AWS_REGION,required"`
	QueueName string `env:"QUEUE_NAME,required"`
}

func main() {
	common.Main(do)
}

func do() error {
	var env env
	if err := appenv.NewManager(allEnvKeys).Populate(&env); err != nil {
		return err
	}
	s := sqs.New(
		&aws.Config{
			Credentials: credentials.NewEnvCredentials(),
			HTTPClient:  cacerts.NewHTTPClient(),
			Region:      env.AwsRegion,
		},
	)
	createQueueOutput, err := s.CreateQueue(
		&sqs.CreateQueueInput{
			QueueName: aws.String(env.QueueName),
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
