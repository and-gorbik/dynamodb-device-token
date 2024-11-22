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

func (r *Repository) GetDeviceList(ctx context.Context, userID int64, kind model.TokenKind) ([]model.Device, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("user_id = :userID and begins_with (kind_device_model, :kind)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID": &types.AttributeValueMemberN{
				Value: strconv.FormatInt(userID, 10),
			},
			":kind": &types.AttributeValueMemberS{
				Value: string(kind),
			},
		},
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	})
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	tokens := make([]model.Device, 0, out.Count)
	for _, item := range out.Items {
		tokens = append(tokens, *FromItem(item))
	}

	log.Printf("[get device list] %s\n", (*printableConsumedCapacity)(out.ConsumedCapacity))
	return tokens, nil
}
