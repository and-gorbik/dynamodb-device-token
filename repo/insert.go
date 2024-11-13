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

func (r *Repository) Put(ctx context.Context, t model.DeviceToken) error {
	out, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			fieldUserID:      &types.AttributeValueMemberN{Value: strconv.FormatInt(t.UserID, 10)},
			fieldKind:        &types.AttributeValueMemberS{Value: string(t.Kind)},
			fieldModifiedAt:  &types.AttributeValueMemberN{Value: strconv.FormatInt(t.ModifiedAt, 10)},
			fieldToken:       &types.AttributeValueMemberS{Value: t.Token},
			fieldAppVersion:  &types.AttributeValueMemberS{Value: t.AppVersion},
			fieldDeviceModel: &types.AttributeValueMemberS{Value: t.DeviceModel},
		},
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	})
	if err != nil {
		return fmt.Errorf("insert one: %w", err)
	}

	log.Printf("[insert one] %s\n", (*printableConsumedCapacity)(out.ConsumedCapacity))
	return nil
}

func (r *Repository) InsertBulk(ctx context.Context, tt []model.DeviceToken) error {
	reqs := make([]types.WriteRequest, 0, len(tt))
	for _, t := range tt {
		reqs = append(reqs, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					fieldUserID:      &types.AttributeValueMemberN{Value: strconv.FormatInt(t.UserID, 10)},
					fieldKind:        &types.AttributeValueMemberS{Value: string(t.Kind)},
					fieldModifiedAt:  &types.AttributeValueMemberN{Value: strconv.FormatInt(t.ModifiedAt, 10)},
					fieldToken:       &types.AttributeValueMemberS{Value: t.Token},
					fieldAppVersion:  &types.AttributeValueMemberS{Value: t.AppVersion},
					fieldDeviceModel: &types.AttributeValueMemberS{Value: t.DeviceModel},
				},
			},
		})
	}

	out, err := r.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			tableName: reqs,
		},
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	})
	if err != nil {
		return fmt.Errorf("batch write item: %w", err)
	}

	for _, req := range out.UnprocessedItems[tableName] {
		log.Printf(
			"item (user_id='%s', kind='%s', token='%s') wasn't processed\n",
			req.PutRequest.Item[fieldUserID],
			req.PutRequest.Item[fieldKind],
			req.PutRequest.Item[fieldToken],
		)
	}

	for _, cc := range out.ConsumedCapacity {
		log.Printf("[insert] %s\n", printableConsumedCapacity(cc))
	}

	return nil
}
