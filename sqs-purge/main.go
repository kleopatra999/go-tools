package main

import (
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/sqs"
	"github.com/peter-edge/go-env"
	"github.com/peter-edge/go-tools/common/aws"
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
