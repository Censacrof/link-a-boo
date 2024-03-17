package shortened_url

import (
	"context"
	"os"

	"github.com/Censacrof/link-a-boo/lambda/go/pkg/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func IsSlugAvailable(ctx context.Context, slug string) (bool, error) {
	dbClient, err := db.GetDbClient(ctx)
	if err != nil {
		return false, err
	}

	key := map[string]types.AttributeValue{
		"slug": &types.AttributeValueMemberS{Value: slug},
	}

	projectionExpression := "slug"
	result, err := dbClient.DdbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:            aws.String(os.Getenv("URLS_TABLE_NAME")),
		Key:                  key,
		ProjectionExpression: &projectionExpression,
	})

	if err != nil {
		return false, err
	}

	if result.Item != nil {
		return false, nil
	}

	return true, nil
}
