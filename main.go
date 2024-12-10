package main

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

func formatDate(timestamps int64) string {
	seconds := timestamps / 1000
	dateTime := time.Unix(seconds, 0)
	formattedDate := dateTime.Format("2006-01-02 15:04:05")
	return formattedDate
}

// func main() {
// 	timestamp := retrieveLogCloudwatch("apiTemplate", 0)
// 	for {
// 		timestamp = retrieveLogCloudwatch("apiTemplate", timestamp) + 1
// 	}
// }

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: stream <command> [options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "login":
		retrieveCreds()

	case "logs":
		stopChan := make(chan struct{})
		lambdaName := os.Args[2]
		timestamp := retrieveLogCloudwatch(lambdaName, 0) + 1
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
				timestamp = retrieveLogCloudwatch(lambdaName, timestamp)
				time.Sleep(1 * time.Second)
			}
		}
	default:
		fmt.Printf("Unknown command %s\n", os.Args[1])
	}
}

func retrieveCreds() {
	dir, err := os.UserHomeDir()
	check(err)
	fmt.Println(dir)

	dat, err := os.ReadFile(dir + "/.aws/credentials")
	if err != nil {
		fmt.Println("Failed to retrieve data ", err)
	}
	fmt.Printf("Here's the content of the file : %s", dat)
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
		StartTime:     aws.Int64(startTime),
		// EndTime:   aws.Int64(currentTime),
	}
	resp, err := client.GetLogEvents(context.TODO(), input)
	// fmt.Println(startTime)
	// seconds := startTime / 1000
	// dateTime := time.Unix(seconds, 0)
	// formattedDate := dateTime.Format("2006-01-02 15:04:05")
	// fmt.Println(formattedDate)
	// fmt.Println(*logStream.LogStreams[0].LogStreamName)
	// fmt.Println(len(resp.Events))
	for _, event := range resp.Events {
		if *event.Timestamp > timestamp {
			timestamp = int64(*event.Timestamp)
		}

		fmt.Printf(*event.Message)
	}
	check(err)

	return timestamp
	// for _, stream := range logStream.LogStreams {
	// 	timestamp := time.Unix(*stream.CreationTime/1000, 0)
	// 	fmt.Println(timestamp.Format("2006-01-02 15:04:05"))
	// 	fmt.Println(*stream.LogStreamName)
	// }
}
