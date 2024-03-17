package shortened_url

import (
	"context"
	"os"

	"github.com/Censacrof/link-a-boo/lambda/go/pkg/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Get(ctx context.Context, slug string) (*shortenedUrl, error) {
	dbClient, err := db.GetDbClient(ctx)
	if err != nil {
		return nil, err
	}

	key := map[string]types.AttributeValue{
		"slug": &types.AttributeValueMemberS{Value: slug},
	}

	item, err := dbClient.DdbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("URLS_TABLE_NAME")),
		Key:       key,
	})

	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	var shortenedUrl shortenedUrl
	err = attributevalue.UnmarshalMap(item.Item, &shortenedUrl)
	if err != nil {
		return nil, err
	}

	return &shortenedUrl, nil
}
