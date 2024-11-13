package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"strconv"

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

	userID, err := strconv.ParseInt(partitionKey, 10, 64)
	if err != nil {
		log.Fatalf("parse int: %v\n", err)
	}

	if sortKey == "" {
		tokens, err := r.GetTokenList(ctx, userID)
		if err != nil {
			log.Fatalf("get token list: %v\n", err)
		}

		printTokens(tokens)
	} else {
		token, err := r.GetToken(ctx, userID, model.TokenKind(sortKey))
		if err != nil {
			log.Fatalf("get token: %v\n", err)
		}

		if token == nil {
			log.Fatalln("token not found")
		}

		printTokens([]model.DeviceToken{*token})
	}
}

func put(ctx context.Context, r *repo.Repository, data string) {
	var token model.DeviceToken
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		log.Fatalf("can't unmarshal item: %v\n", err)
	}

	if err := r.Put(ctx, token); err != nil {
		log.Fatalf("put: %v\n", err)
	}

	log.Println("item is put successfully")
}

func delete(ctx context.Context, r *repo.Repository, partitionKey, sortKey string) {
	if partitionKey == "" {
		log.Fatalln("partition key is required")
	}

	userID, err := strconv.ParseInt(partitionKey, 10, 64)
	if err != nil {
		log.Fatalf("parse int: %v\n", err)
	}

	if sortKey == "" {
		if err := r.DeleteTokens(ctx, userID); err != nil {
			log.Fatalf("delete tokens: %v\n", err)
		}
	} else {
		token, err := r.DeleteToken(ctx, userID, model.TokenKind(sortKey))
		if err != nil {
			log.Fatalf("delete token: %v\n", err)
		}

		if token == nil {
			log.Fatalln("token not found")
		}

		printTokens([]model.DeviceToken{*token})
	}
}

func printTokens(tt []model.DeviceToken) {
	for _, t := range tt {
		log.Printf("{userID: %d, token: %s, kind: %s}\n", t.UserID, t.Token, t.Kind)
	}
}
