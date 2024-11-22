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

func (r *Repository) DeleteToken(ctx context.Context, userID int64, kind model.TokenKind, deviceModel string) (*model.Device, error) {
	d := &model.Device{
		UserID:      userID,
		Kind:        kind,
		DeviceModel: deviceModel,
	}

	// нельзя удалять больше одного элемента за раз
	out, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName:              aws.String(tableName),
		Key:                    ToKey(d.PartitionKey(), d.SortKey()),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityIndexes,
		ReturnValues:           types.ReturnValueAllOld,
	})
	if err != nil {
		return nil, fmt.Errorf("delete item: %w", err)
	}

	log.Printf("[delete token] %s\n", (*printableConsumedCapacity)(out.ConsumedCapacity))

	if len(out.Attributes) == 0 {
		return nil, nil
	}

	return FromItem(out.Attributes), nil
}
