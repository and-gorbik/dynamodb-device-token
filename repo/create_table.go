package repo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (r *Repository) CreateTable(ctx context.Context) (*types.TableDescription, error) {
	table, err := r.client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(fieldUserID),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String(fieldKindDeviceModel),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},

		// if sort key is defined, partition + sort keys must be unique together
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(fieldUserID),
				KeyType:       types.KeyTypeHash, // partition key - required, exactly one
			},
			{
				AttributeName: aws.String(fieldKindDeviceModel),
				KeyType:       types.KeyTypeRange, // sort key - not required, only one is possibly
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	// если выйти сейчас, вернется описание еще не созданной таблицы
	// return table.TableDescription, nil
	_ = table

	// Инициализируем waiter и ждем создания таблицы
	waiter := dynamodb.NewTableExistsWaiter(r.client)
	err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}, maxDurationOfTableCreation)
	if err != nil {
		return nil, fmt.Errorf("error waiting for table to become active: %w", err)
	}

	fmt.Println("Table is now active!")

	if err := r.enableTTL(ctx); err != nil {
		return nil, fmt.Errorf("enable ttl: %w", err)
	}

	out, err := r.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: &tableName,
	})
	if err != nil {
		return nil, fmt.Errorf("describe table: %w", err)
	}

	return out.Table, nil
}

func (r *Repository) enableTTL(ctx context.Context) error {
	_, err := r.client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String(tableName),
		TimeToLiveSpecification: &types.TimeToLiveSpecification{
			AttributeName: aws.String(fieldTTL),
			Enabled:       aws.Bool(true),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
