package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/and-gorbik/dynamodb-device-token/db"
	"github.com/and-gorbik/dynamodb-device-token/model"
	"github.com/and-gorbik/dynamodb-device-token/region"
	"github.com/and-gorbik/dynamodb-device-token/repo"
)

func main() {
	command := flag.String("command", "", "command to apply; can either be 'create', 'apply', 'delete', 'enable-stream' or 'make-global-table'")
	reg := flag.String("region", region.Default(), "set up the region, from where global table will be created")
	replicaRegion := flag.String("replica-region", "", "region where replicas will be created")
	flag.Parse()

	if *command == "" {
		log.Fatal("param 'command' is required")
	}

	if !region.In(*reg) {
		log.Fatalf("unknown region: %s\n", *reg)
	}

	ctx := context.Background()

	dev, _ := strconv.ParseBool(os.Getenv("DEV"))

	client := db.InitDynamoDBClient(ctx, *reg, dev)
	r := repo.Init(client)

	switch *command {
	case "create":
		create(ctx, r)
	case "apply":
		apply(ctx, r)
	case "delete":
		delete(ctx, r)
	case "enable-stream":
		enableStream(ctx, r)
	case "make-global-table":
		if !region.In(*replicaRegion) {
			log.Fatalf("unknown replica region %s\n", *replicaRegion)
		}

		makeGlobalTable(ctx, r, *replicaRegion)
	default:
		log.Fatalf("unknown command: %s\n", *command)
	}
}

func create(ctx context.Context, r *repo.Repository) {
	desc, err := r.CreateTable(ctx)
	if err != nil {
		log.Fatalf("create table: %s\n", err)
	}

	printTableDescription(desc)
}

func apply(ctx context.Context, r *repo.Repository) {
	tokens, err := parseRecords("./records.json")
	if err != nil {
		log.Fatalf("apply: %s\n", err)
	}

	if err := r.InsertBulk(ctx, tokens); err != nil {
		log.Fatalf("insert bulk: %s\n", err)
	}
}

func delete(ctx context.Context, r *repo.Repository) {
	if err := r.DropTable(ctx); err != nil {
		log.Fatalf("drop table: %s\n", err)
	}
}

func enableStream(ctx context.Context, r *repo.Repository) {
	streamID, err := r.EnableStreaming(ctx)
	if err != nil {
		log.Fatalf("enable stream: %s\n", err)
	}

	log.Println("stream id: ", streamID)
}

func makeGlobalTable(ctx context.Context, r *repo.Repository, replicaRegion string) {
	if err := r.MakeGlobalTable(ctx, replicaRegion); err != nil {
		log.Fatalf("make global table: %v\n", err)
	}

	log.Printf("replica is created in region %s\n", replicaRegion)
}

func parseRecords(path string) ([]model.Device, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("parse records: %w", err)
	}

	var dd []model.Device
	if err := json.Unmarshal(data, &dd); err != nil {
		return nil, fmt.Errorf("unmarshal records: %w", err)
	}

	return dd, nil
}

func printTableDescription(desc *types.TableDescription) {
	log.Println("=== TABLE DESCRIPTION ===")
	log.Printf("Name: %s\n", *desc.TableName)
	log.Printf("Status: %s\n", desc.TableStatus)
	log.Printf("Attributes:\n")
	for _, attr := range desc.AttributeDefinitions {
		log.Printf("* %s [%s]\n", *attr.AttributeName, attr.AttributeType)
	}

	log.Printf("Key schema:\n")
	for _, elem := range desc.KeySchema {
		log.Printf("* %s [%s]", *elem.AttributeName, elem.KeyType)
	}
}
