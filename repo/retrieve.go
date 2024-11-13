package repo

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/and-gorbik/dynamodb-device-token/model"
)

func (r *Repository) GetTokenList(ctx context.Context, userID int64) ([]model.DeviceToken, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("user_id = :userID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID": &types.AttributeValueMemberN{
				Value: strconv.FormatInt(userID, 10),
			},
		},
		ConsistentRead:         aws.Bool(false),
		ScanIndexForward:       aws.Bool(true),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityIndexes,
	})
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	tokens := make([]model.DeviceToken, 0, out.Count)
	for _, item := range out.Items {
		modifiedAt, _ := strconv.ParseInt(item[fieldModifiedAt].(*types.AttributeValueMemberN).Value, 10, 64)

		tokens = append(tokens, model.DeviceToken{
			UserID:      userID,
			Kind:        model.TokenKind(item[fieldKind].(*types.AttributeValueMemberS).Value),
			Token:       item[fieldToken].(*types.AttributeValueMemberS).Value,
			DeviceModel: item[fieldDeviceModel].(*types.AttributeValueMemberS).Value,
			AppVersion:  item[fieldAppVersion].(*types.AttributeValueMemberS).Value,
			ModifiedAt:  modifiedAt,
		})
	}

	log.Printf("[get token list] %s\n", (*printableConsumedCapacity)(out.ConsumedCapacity))
	return tokens, nil
}

func (r *Repository) GetToken(ctx context.Context, userID int64, kind model.TokenKind) (*model.DeviceToken, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberN{
				Value: strconv.FormatInt(userID, 10),
			},
			"kind": &types.AttributeValueMemberS{
				Value: string(kind),
			},
		},

		ConsistentRead:         aws.Bool(false),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityIndexes,
	})
	if err != nil {
		return nil, fmt.Errorf("get item: %w", err)
	}

	if len(out.Item) == 0 {
		return nil, nil
	}

	modifiedAt, _ := strconv.ParseInt(out.Item[fieldModifiedAt].(*types.AttributeValueMemberN).Value, 10, 64)

	log.Printf("[get token] %s\n", (*printableConsumedCapacity)(out.ConsumedCapacity))
	return &model.DeviceToken{
		UserID:      userID,
		Kind:        model.TokenKind(out.Item[fieldKind].(*types.AttributeValueMemberS).Value),
		Token:       out.Item[fieldToken].(*types.AttributeValueMemberS).Value,
		DeviceModel: out.Item[fieldDeviceModel].(*types.AttributeValueMemberS).Value,
		AppVersion:  out.Item[fieldAppVersion].(*types.AttributeValueMemberS).Value,
		ModifiedAt:  modifiedAt,
	}, nil
}
