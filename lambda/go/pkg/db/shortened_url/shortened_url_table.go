package shortened_url

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/Censacrof/link-a-boo/lambda/go/pkg/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

func Put(ctx context.Context, shortenedUrl ShortenedUrl) error {
	dbClient, err := db.GetDbClient(ctx)
	if err != nil {
		return err
	}

	item, err := attributevalue.MarshalMap(shortenedUrl)

	tableName := os.Getenv("URLS_TABLE_NAME")

	if err != nil {
		return fmt.Errorf("Put in table '%s' failed: %w", tableName, err)
	}

	_, err = dbClient.DdbClient.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	})

	if err != nil {
		return fmt.Errorf("Put in table '%s' failed: %w", tableName, err)
	}

	return nil
}

func Get(ctx context.Context, slug string) (*ShortenedUrl, error) {
	dbClient, err := db.GetDbClient(ctx)
	if err != nil {
		return nil, err
	}

	key := map[string]types.AttributeValue{
		slug: &types.AttributeValueMemberS{Value: slug},
	}

	item, err := dbClient.DdbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key: key,
	})

	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	var shortenedUrl ShortenedUrl
	err = attributevalue.UnmarshalMap(item.Item, shortenedUrl)
	if err != nil {
		return nil, err
	}

	return &shortenedUrl, nil
}

type ShortenedUrl struct {
	Slug string `dynamodbav:"slug"`
	Url  string `dynamodbav:"url"`
}

func NewShortenedUrl(url url.URL) ShortenedUrl {
	slug := uuid.New().String()

	return ShortenedUrl{
		Slug: slug,
		Url:  url.String(),
	}
}
