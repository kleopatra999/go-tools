package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/codeship/go-tools/common/aws"
	"github.com/peter-edge/go-env"
)

type appEnv struct {
	QueueName string `env:"QUEUE_NAME,required"`
}

func main() {
	commonaws.Main(do)
}

func do(awsConfig *aws.Config) error {
	var appEnv appEnv
	if err := env.Populate(&appEnv, env.PopulateOptions{}); err != nil {
		return err
	}
	sqsClient := sqs.New(awsConfig)
	createQueueOutput, err := sqsClient.CreateQueue(
		&sqs.CreateQueueInput{
			QueueName: aws.String(appEnv.QueueName),
		},
	)
	if err != nil {
		return err
	}
	_, err = sqsClient.PurgeQueue(
		&sqs.PurgeQueueInput{
			QueueURL: createQueueOutput.QueueURL,
		},
	)
	return err
}
