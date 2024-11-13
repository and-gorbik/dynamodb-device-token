package repo

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	tableName = "device_token"

	fieldUserID      = "user_id"
	fieldKind        = "kind"
	fieldModifiedAt  = "modified_at"
	fieldToken       = "token"
	fieldAppVersion  = "app_version"
	fieldDeviceModel = "device_model"

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
