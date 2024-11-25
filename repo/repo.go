package repo

import (
	"fmt"
	"strconv"
	"time"

	"github.com/and-gorbik/dynamodb-device-token/model"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	tableName = "device"

	fieldUserID          = "user_id"
	fieldKindDeviceModel = "kind_device_model"
	fieldModifiedAt      = "modified_at"
	fieldToken           = "token"
	fieldAppVersion      = "app_version"
	fieldLocale          = "locale"
	fieldTTL             = "ttl"

	maxDurationOfTableCreation = time.Minute * 5
)

type Repository struct {
	client *dynamodb.Client
}

func Init(c *dynamodb.Client) *Repository {
	return &Repository{c}
}

type printableConsumedCapacity types.ConsumedCapacity

func (p printableConsumedCapacity) String() string {
	if p.Table == nil {
		p.Table = &types.Capacity{}
	}

	return fmt.Sprintf(
		"total: %s, rcu: %s, wcu: %s, table rcu: %s, table wcu: %s\n",
		printTPtr(p.CapacityUnits),
		printTPtr(p.ReadCapacityUnits),
		printTPtr(p.WriteCapacityUnits),
		printTPtr(p.Table.ReadCapacityUnits),
		printTPtr(p.Table.WriteCapacityUnits),
	)
}

func printTPtr[T any](val *T) string {
	if val == nil {
		return "-"
	}

	return fmt.Sprintf("%v", *val)
}

func ToItem(d *model.Device) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		fieldUserID:          &types.AttributeValueMemberN{Value: d.PartitionKey()},
		fieldKindDeviceModel: &types.AttributeValueMemberS{Value: d.SortKey()},
		fieldModifiedAt:      &types.AttributeValueMemberN{Value: strconv.FormatInt(d.ModifiedAt, 10)},
		fieldToken:           &types.AttributeValueMemberS{Value: d.Token},
		fieldAppVersion:      &types.AttributeValueMemberS{Value: d.AppVersion},
		fieldLocale:          &types.AttributeValueMemberS{Value: d.Locale},
		fieldTTL:             &types.AttributeValueMemberN{Value: strconv.FormatInt(d.TTL, 10)},
	}
}

func ToKey(pk, sk string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		fieldUserID:          &types.AttributeValueMemberN{Value: pk},
		fieldKindDeviceModel: &types.AttributeValueMemberS{Value: sk},
	}
}

func ToPartitionKey(pk string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		fieldUserID: &types.AttributeValueMemberN{Value: pk},
	}
}

func FromItem(item map[string]types.AttributeValue) *model.Device {
	userID, _ := strconv.ParseInt(item[fieldUserID].(*types.AttributeValueMemberN).Value, 10, 64)
	modifiedAt, _ := strconv.ParseInt(item[fieldModifiedAt].(*types.AttributeValueMemberN).Value, 10, 64)

	d := model.Device{
		UserID:     userID,
		ModifiedAt: modifiedAt,
		Token:      item[fieldToken].(*types.AttributeValueMemberS).Value,
		AppVersion: item[fieldAppVersion].(*types.AttributeValueMemberS).Value,
		Locale:     item[fieldLocale].(*types.AttributeValueMemberS).Value,
	}

	d.SetSortKey(item[fieldKindDeviceModel].(*types.AttributeValueMemberS).Value)

	return &d
}
