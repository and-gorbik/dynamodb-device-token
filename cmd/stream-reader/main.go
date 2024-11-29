package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"

	"github.com/and-gorbik/dynamodb-device-token/db"
	"github.com/and-gorbik/dynamodb-device-token/region"
)

func main() {
	streamID := flag.String("stream-id", "", "stream id from './cmd/db-manager --command enable-stream'")
	reg := flag.String("region", region.Default(), "set up the region")
	flag.Parse()

	ctx := context.Background()

	if *streamID == "" {
		log.Fatal("stream id is not provided")
	}

	if !region.In(*reg) {
		log.Fatalf("unknown region: %s\n", *reg)
	}

	client := db.InitDynamoDBStreamClient(ctx, *reg)
	startStreaming(ctx, client, *streamID)
}

func startStreaming(ctx context.Context, client *dynamodbstreams.Client, streamID string) {
	describeStreamOutput, err := client.DescribeStream(ctx, &dynamodbstreams.DescribeStreamInput{
		StreamArn: &streamID,
	})
	if err != nil {
		log.Fatalf("failed to describe stream: %v\n", err)
	}

	if len(describeStreamOutput.StreamDescription.Shards) == 0 {
		log.Fatalln("no available shards")
	}

	shard := describeStreamOutput.StreamDescription.Shards[0]

	// Получаем итератор для чтения данных с начала потока
	shardIteratorOutput, err := client.GetShardIterator(ctx, &dynamodbstreams.GetShardIteratorInput{
		StreamArn:         &streamID,
		ShardId:           shard.ShardId,
		ShardIteratorType: types.ShardIteratorTypeLatest,
	})
	if err != nil {
		log.Fatalf("failed to get shard iterato: %v\n", err)
	}

	shardIterator := shardIteratorOutput.ShardIterator

	// Непрерывное чтение данных из потока
	for {
		if shardIterator == nil {
			log.Println("Shard iterator expired, exit")
			break
		}

		// Получаем записи из потока
		getRecordsOutput, err := client.GetRecords(ctx, &dynamodbstreams.GetRecordsInput{
			ShardIterator: shardIterator,
		})
		if err != nil {
			log.Fatalf("failed to get records, %v", err)
		}

		// Обрабатываем каждую запись
		for _, record := range getRecordsOutput.Records {
			fmt.Printf("Event: %s\n", record.EventName)
			fmt.Printf("Keys: %v\n", record.Dynamodb.Keys)
			if record.Dynamodb.NewImage != nil {
				fmt.Printf("New Image: %v\n", record.Dynamodb.NewImage)
			}
			if record.Dynamodb.OldImage != nil {
				fmt.Printf("Old Image: %v\n", record.Dynamodb.OldImage)
			}
		}

		// Обновляем итератор для следующего цикла
		shardIterator = getRecordsOutput.NextShardIterator

		// Задержка, чтобы избежать превышения лимита запросов к DynamoDB Streams
		time.Sleep(2 * time.Second)
	}
}
