package main

import (
	"fmt"
	"sort"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ec2"
	"github.com/awslabs/aws-sdk-go/service/ecs"
	"github.com/peter-edge/go-env"
	"github.com/peter-edge/go-tools/common/aws"
)

type appEnv struct {
	EcsCluster string `env:"ECS_CLUSTER,required"`
}

func main() {
	commonaws.Main(do)
}

func do(awsConfig *aws.Config) error {
	var appEnv appEnv
	if err := env.Populate(&appEnv, env.PopulateOptions{}); err != nil {
		return err
	}
	ecsClient := ecs.New(awsConfig)
	listContainerInstancesOutput, err := ecsClient.ListContainerInstances(
		&ecs.ListContainerInstancesInput{
			Cluster: aws.String(appEnv.EcsCluster),
		},
	)
	if err != nil {
		return err
	}
	describeContainerInstancesOutput, err := ecsClient.DescribeContainerInstances(
		&ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String(appEnv.EcsCluster),
			ContainerInstances: listContainerInstancesOutput.ContainerInstanceARNs,
		},
	)
	if err != nil {
		return err
	}
	ec2InstanceIDs := make([]*string, len(describeContainerInstancesOutput.ContainerInstances))
	for i, containerInstance := range describeContainerInstancesOutput.ContainerInstances {
		ec2InstanceIDs[i] = containerInstance.EC2InstanceID
	}
	ec2Client := ec2.New(awsConfig)
	describeInstancesOutput, err := ec2Client.DescribeInstances(
		&ec2.DescribeInstancesInput{
			InstanceIDs: ec2InstanceIDs,
		},
	)
	if err != nil {
		return err
	}
	var publicDNSNames []string
	for _, reservation := range describeInstancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			if instance.PublicDNSName != nil {
				publicDNSNames = append(publicDNSNames, *instance.PublicDNSName)
			}
		}
	}
	sort.Sort(sort.StringSlice(publicDNSNames))
	for i, publicDNSName := range publicDNSNames {
		fmt.Printf("%d: %s\n", i, publicDNSName)
	}
	return nil
}
