package db

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"
	"github.com/joho/godotenv"
)

func InitDynamoDBStreamClient(ctx context.Context, region string) *dynamodbstreams.Client {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return dynamodbstreams.NewFromConfig(cfg)
}
