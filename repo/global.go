package repo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (r *Repository) MakeGlobalTable(ctx context.Context, region string) error {
	_, err := r.client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
		TableName: aws.String(tableName),
		ReplicaUpdates: []types.ReplicationGroupUpdate{
			{
				Create: &types.CreateReplicationGroupMemberAction{
					RegionName: aws.String(region),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("update table: %w", err)
	}

	return nil
}
