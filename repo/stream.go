package repo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (r *Repository) EnableStreaming(ctx context.Context) (string, error) {
	out, err := r.client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
		TableName: aws.String(tableName),
		StreamSpecification: &types.StreamSpecification{
			StreamEnabled:  aws.Bool(true),
			StreamViewType: types.StreamViewTypeNewAndOldImages,
		},
	})
	if err != nil {
		return "", fmt.Errorf("enable stream for the table: %w", err)
	}

	streamID := out.TableDescription.LatestStreamArn
	if streamID == nil {
		return "", fmt.Errorf("stream arn is nil")
	}

	return *streamID, nil
}
