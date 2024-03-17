package shortened_url

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/Censacrof/link-a-boo/lambda/go/pkg/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const UrlMaxLength int = 2048
const SlugMaxLength int = 50
const SlugMinLength int = 4

var ErrSlugAlreadyExists = errors.New("Provided slug already exists")

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

type shortenedUrl struct {
	Slug string  `dynamodbav:"slug"`
	Url  url.URL `dynamodbav:"url"`
}

var ErrInvalidSlug = errors.New("Slug is not valid")
var ErrInvalidUrl = errors.New("Url is not valid")

func New(rawUrl string, slug string) (*shortenedUrl, error) {
	if len(slug) > SlugMaxLength {
		return nil, ErrInvalidSlug
	}

	if len(slug) < SlugMinLength {
		return nil, ErrInvalidSlug
	}

	if len(rawUrl) > UrlMaxLength {
		return nil, ErrInvalidUrl
	}

	url, err := url.Parse(rawUrl)
	if err != nil {
		return nil, ErrInvalidUrl
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return nil, ErrInvalidUrl
	}

	if url.Host == "" {
		return nil, ErrInvalidUrl
	}

	return &shortenedUrl{
		Slug: slug,
		Url:  *url,
	}, nil
}
