package instances

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/sirupsen/logrus"
	"log"
	"sort"
	"time"
)

var instances map[string]types.Instance

func init() {
	logrus.SetLevel(logrus.WarnLevel)
	instances = make(map[string]types.Instance)
}

func InstanceDate(cfg aws.Config, instance string) time.Time {
	instanceType := GetInstance(cfg, instance)
	return *instanceType.LaunchTime
}

func InstanceName(cfg aws.Config, instance string) string {
	instanceType := GetInstance(cfg, instance)
	for _, tag := range instanceType.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return "Not Found"
}

func Instances(cfg aws.Config) (instanceIds []string) {
	client := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeInstancesInput{}
	output, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		log.Fatal(err)
	}
	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			setInstance(instance)
			instanceIds = append(instanceIds, *instance.InstanceId)
		}
	}
	sort.Strings(instanceIds)
	return instanceIds
}

func setInstance(instance types.Instance) {
	logrus.Trace("Setting instance: " + *instance.InstanceId)
	instances[*instance.InstanceId] = instance
}

func GetInstance(cfg aws.Config, instanceId string) types.Instance {
	logrus.Trace("Getting instance: " + instanceId)
	thisInstance, exists := instances[instanceId]
	if exists {
		return thisInstance
	} else {
		client := ec2.NewFromConfig(cfg)
		input := &ec2.DescribeInstancesInput{
			InstanceIds: []string{instanceId},
		}
		output, err := client.DescribeInstances(context.TODO(), input)
		if err != nil {
			log.Fatal(err)
		}
		for _, reservation := range output.Reservations {
			for _, instance := range reservation.Instances {
				setInstance(instance)
				return instance
			}
		}
	}
	return types.Instance{}
}
