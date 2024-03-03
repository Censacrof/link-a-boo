package db

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

type ShortenedUrlTable struct {
}

var shortenedUrlTable *ShortenedUrlTable = nil

func GetShortenedUrlTable() *ShortenedUrlTable {
	if shortenedUrlTable == nil {
		shortenedUrlTable = &ShortenedUrlTable{}
	}

	return shortenedUrlTable
}

func (self *ShortenedUrlTable) Put(ctx context.Context, shortenedUrl ShortenedUrl) error {
	dbClient, err := GetDbClient(ctx)
	if err != nil {
		return err
	}

	item, err := attributevalue.MarshalMap(shortenedUrl)

	tableName := os.Getenv("URLS_TABLE_NAME")

	if err != nil {
		return fmt.Errorf("Put in table '%s' failed: %w", tableName, err)
	}

	_, err = dbClient.ddbClient.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	})

	if err != nil {
		return fmt.Errorf("Put in table '%s' failed: %w", tableName, err)
	}

	return nil
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
