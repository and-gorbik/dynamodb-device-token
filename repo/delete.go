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

func (r *Repository) DeleteToken(ctx context.Context, userID int64, kind model.TokenKind) (*model.DeviceToken, error) {
	key := map[string]types.AttributeValue{
		"user_id": &types.AttributeValueMemberN{
			Value: strconv.FormatInt(userID, 10),
		},
		"kind": &types.AttributeValueMemberS{
			Value: string(kind),
		},
	}

	// нельзя удалять больше одного элемента за раз
	out, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName:              aws.String(tableName),
		Key:                    key,
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

	modifiedAt, _ := strconv.ParseInt(out.Attributes[fieldModifiedAt].(*types.AttributeValueMemberN).Value, 10, 64)

	return &model.DeviceToken{
		UserID:      userID,
		Kind:        model.TokenKind(out.Attributes[fieldKind].(*types.AttributeValueMemberS).Value),
		Token:       out.Attributes[fieldToken].(*types.AttributeValueMemberS).Value,
		DeviceModel: out.Attributes[fieldDeviceModel].(*types.AttributeValueMemberS).Value,
		AppVersion:  out.Attributes[fieldAppVersion].(*types.AttributeValueMemberS).Value,
		ModifiedAt:  modifiedAt,
	}, nil
}

func (r *Repository) DeleteTokens(ctx context.Context, userID int64) error {
	reqs := make([]types.WriteRequest, 0, 4)
	for _, kind := range []model.TokenKind{model.TokenKindAndroidGeneral, model.TokenKindIOSGeneral, model.TokenKindIOSLiveActivity, model.TokenKindIOSVoip} {
		reqs = append(reqs, types.WriteRequest{

			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"user_id": &types.AttributeValueMemberN{
						Value: strconv.FormatInt(userID, 10),
					},
					"kind": &types.AttributeValueMemberS{
						Value: string(kind),
					},
				},
			},
		})
	}

	// удалять элементы можно только явно указав partition key + sort key,
	// поэтому приходится явно задавать батч из всех записей на удаление.
	// такой батч не может быть больше 25 элементов.
	out, err := r.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			tableName: reqs,
		},
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	})
	if err != nil {
		return fmt.Errorf("delete tokens: %w", err)
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
		log.Printf("[delete tokens] %s\n", printableConsumedCapacity(cc))
	}

	return nil
}
