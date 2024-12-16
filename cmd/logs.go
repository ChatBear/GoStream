package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func retrieveLogCloudwatch(name string, timestamp int64) int64 {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	check(err)

	client := cloudwatchlogs.NewFromConfig(cfg)
	logGroupName := fmt.Sprintf("/aws/lambda/%s", name)

	inputStream := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroupName),
		Limit:        aws.Int32(1),
		Descending:   aws.Bool(true),
		OrderBy:      "LastEventTime",
	}

	logStream, err := client.DescribeLogStreams(context.TODO(), inputStream)
	check(err)

	startTime := timestamp
	input := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(*logStream.LogStreams[0].LogStreamName),
		StartTime:     aws.Int64(startTime + 1),
	}
	resp, err := client.GetLogEvents(context.TODO(), input)

	for _, event := range resp.Events {
		if *event.Timestamp > timestamp {
			timestamp = int64(*event.Timestamp)
		}

		fmt.Printf(*event.Message)
	}
	check(err)

	return timestamp
}

func BeginLogStream(logName string) {
	stopChan := make(chan struct{})
	timestamp := retrieveLogCloudwatch(logName, 0) + 1
	go func() {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Press enter to stop")
		_, _ = reader.ReadByte()
		close(stopChan)
	}()
	for {
		select {
		case <-stopChan:
			fmt.Println("Loop stopped by user.")
			return
		default:
			timestamp = retrieveLogCloudwatch(logName, timestamp)
			time.Sleep(1 * time.Second)
		}
	}
}
