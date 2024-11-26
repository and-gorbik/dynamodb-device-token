package repo

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/and-gorbik/dynamodb-device-token/model"
)

func (r *Repository) Put(ctx context.Context, d model.Device) error {
	d.SetTTL()

	if !d.Latest {
		out, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName:              aws.String(tableName),
			Item:                   ToItem(&d, false),
			ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		})
		if err != nil {
			return fmt.Errorf("put item: %w", err)
		}

		log.Printf("[insert one] %s\n", (*printableConsumedCapacity)(out.ConsumedCapacity))
		return nil
	}

	out, err := r.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			tableName: {
				{
					PutRequest: &types.PutRequest{
						Item: ToItem(&d, true),
					},
				},
				{
					PutRequest: &types.PutRequest{
						Item: ToItem(&d, false),
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("insert two records: %w", err)
	}

	for _, cc := range out.ConsumedCapacity {
		log.Printf("[insert] %s\n", printableConsumedCapacity(cc))
	}

	return nil
}

func (r *Repository) InsertBulk(ctx context.Context, dd []model.Device) error {
	reqs := make([]types.WriteRequest, 0, len(dd))
	for _, d := range dd {
		d.SetTTL()
		reqs = append(reqs, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: ToItem(&d, false),
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
			"item (user_id='%s', kind_device_model='%s', token='%s') wasn't processed\n",
			req.PutRequest.Item[fieldUserID],
			req.PutRequest.Item[fieldKindDeviceModel],
			req.PutRequest.Item[fieldToken],
		)
	}

	for _, cc := range out.ConsumedCapacity {
		log.Printf("[insert] %s\n", printableConsumedCapacity(cc))
	}

	return nil
}
