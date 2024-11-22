package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/and-gorbik/dynamodb-device-token/db"
	"github.com/and-gorbik/dynamodb-device-token/model"
	"github.com/and-gorbik/dynamodb-device-token/repo"
)

func main() {
	command := flag.String("command", "", "command to apply; can either be 'get', 'put', 'delete'")
	pk := flag.String("pk", "", "partition key value")
	sk := flag.String("sort", "", "sort key value")
	data := flag.String("data", "", "item represented as json")
	flag.Parse()

	if *command == "" {
		log.Fatal("param 'command' is required")
	}

	ctx := context.Background()

	client := db.InitDynamoDBClient(ctx)
	r := repo.Init(client)

	switch *command {
	case "get":
		get(ctx, r, *pk, *sk)
	case "put":
		put(ctx, r, *data)
	case "delete":
		delete(ctx, r, *pk, *sk)
	}
}

func get(ctx context.Context, r *repo.Repository, partitionKey, sortKey string) {
	if partitionKey == "" {
		log.Fatalln("partition key is required")
	}

	if sortKey == "" {
		log.Fatalln("sort key is required")
	}

	userID, err := strconv.ParseInt(partitionKey, 10, 64)
	if err != nil {
		log.Fatalf("parse int: %v\n", err)
	}

	devices, err := r.GetDeviceList(ctx, userID, model.TokenKind(sortKey))
	if err != nil {
		log.Fatalf("get token list: %v\n", err)
	}

	printTokens(devices)
}

func put(ctx context.Context, r *repo.Repository, data string) {
	var d model.Device
	if err := json.Unmarshal([]byte(data), &d); err != nil {
		log.Fatalf("can't unmarshal item: %v\n", err)
	}

	if err := r.Put(ctx, d); err != nil {
		log.Fatalf("put: %v\n", err)
	}

	log.Println("item is put successfully")
}

func delete(ctx context.Context, r *repo.Repository, partitionKey, sortKey string) {
	if partitionKey == "" {
		log.Fatalln("partition key is required")
	}

	if sortKey == "" {
		log.Fatalln("sort key is required")
	}

	userID, err := strconv.ParseInt(partitionKey, 10, 64)
	if err != nil {
		log.Fatalf("parse int: %v\n", err)
	}

	parts := strings.Split(sortKey, "#")
	if len(parts) != 2 {
		log.Fatalf("invalid sort key: it must contain '#'")
	}

	token, err := r.DeleteToken(ctx, userID, model.TokenKind(parts[0]), parts[1])
	if err != nil {
		log.Fatalf("delete token: %v\n", err)
	}

	if token == nil {
		log.Fatalln("token not found")
	}

	printTokens([]model.Device{*token})
}

func printTokens(tt []model.Device) {
	for _, t := range tt {
		log.Printf("{userID: %d, token: %s, kind: %s, device_model: %s}\n", t.UserID, t.Token, t.Kind, t.DeviceModel)
	}
}
