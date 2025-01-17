package db

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
)

/*
Next env variables is needed to be set up in the .env file:
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
*/
func InitDynamoDBClient(ctx context.Context, region string, dev bool) *dynamodb.Client {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	opts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	if dev {
		opts = append(opts, config.WithBaseEndpoint("http://localhost:8000"))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return dynamodb.NewFromConfig(cfg)
}
