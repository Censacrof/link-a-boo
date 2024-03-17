package shortened_url

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Censacrof/link-a-boo/lambda/go/pkg/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Put(ctx context.Context, shortenedUrl *shortenedUrl) error {
	dbClient, err := db.GetDbClient(ctx)
	if err != nil {
		return err
	}

	item, err := attributevalue.MarshalMap(shortenedUrl)

	tableName := os.Getenv("URLS_TABLE_NAME")

	if err != nil {
		return fmt.Errorf("Put in table '%s' failed: %w", tableName, err)
	}

	condition := "attribute_not_exists(slug)"
	_, err = dbClient.DdbClient.PutItem(ctx, &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String(tableName),
		ConditionExpression: &condition,
	})

	if err != nil {
		if e := new(types.ConditionalCheckFailedException); errors.As(err, &e) {
			return ErrSlugAlreadyExists
		}
		return fmt.Errorf("Put in table '%s' failed: %w", tableName, err)
	}

	return nil
}
