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
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: my-cli-app <command> [options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "login":
		retrieveCreds()

	case "logs":
		stopChan := make(chan struct{})

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
				retrieveLogCloudwatch("apiTemplate")
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

func retrieveLogCloudwatch(name string) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	check(err)

	client := cloudwatchlogs.NewFromConfig(cfg)
	logGroupName := fmt.Sprintf("/aws/lambda/%s", name)

	// currentTime := time.Now().Unix() * 1000 // Convert to milliseconds
	// startTime := currentTime - 1000
	// input := &cloudwatchlogs.GetLogEventsInput{
	// 	LogGroupName: aws.String(logGroupName),
	// 	// Limit:        aws.Int32(10),
	// 	// StartTime: aws.Int64(startTime),
	// 	// EndTime:      aws.Int64(currentTime),

	// }

	// fmt.Println(input)
	inputStream := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroupName),
		Descending:   aws.Bool(true),
	}

	// resp, err := client.GetLogEvents(context.TODO(), input)
	logStream, err := client.DescribeLogStreams(context.TODO(), inputStream)
	check(err)

	// if len(resp.Events) == 0 {
	// 	fmt.Printf("\nNo Events\n")
	// 	return
	// }
	fmt.Println(logStream.LogStreams)

	for _, stream := range logStream.LogStreams {
		timestamp := time.Unix(*stream.CreationTime/1000, 0)
		fmt.Println(timestamp.Format("2006-01-02 15:04:05"))

	}
	// for _, event := range resp.Events {
	// 	// timestamp := time.Unix(*event.Timestamp/1000, 0)
	// 	// // Format the timestamp to a classic date format
	// 	// fmt.Println(timestamp.Format("2006-01-02 15:04:05"))
	// 	fmt.Printf(*event.Message)
	// }
	// for event_index_desc_order := len(resp.Events) - 1; event_index_desc_order >= 0; event_index_desc_order-- {
	// 	fmt.Printf(*resp.Events[event_index_desc_order].Message)
	// }
}
