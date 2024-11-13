package db

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
)

const (
	region = "us-east-1"
)

/*
Next env variables is needed to be set up in the .env file:
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
*/
func InitDynamoDBClient(ctx context.Context) *dynamodb.Client {
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

	return dynamodb.NewFromConfig(cfg)
}
