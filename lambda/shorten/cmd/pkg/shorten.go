package shorten

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

type ShortenRequest struct {
	TargetUrl string `json:"targetUrl"`
}

type ShortenResponse struct {
	ShortenedUrl string `json:"shortenedUrl"`
}

func HandleShortenRequest(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Can't load default configuration: %v", err),
		}, nil

	}

	var shortenRequest ShortenRequest
	err = json.Unmarshal([]byte(event.Body), &shortenRequest)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request",
		}, nil
	}

	targetUrl, err := url.Parse(shortenRequest.TargetUrl)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid targetUrl",
		}, nil
	}

	ddbClient := dynamodb.NewFromConfig(cfg)
	shortenedUrl, err := AddUrlToTable(ctx, *ddbClient, *targetUrl)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Internal server error: %v", err),
		}, nil
	}

	resp := ShortenResponse{
		ShortenedUrl: shortenedUrl,
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(respBody),
	}, nil
}

type ShortenedUrlDocument struct {
	Slug string `dynamodbav:"slug"`
	Url  string `dynamodbav:"url"`
}

func AddUrlToTable(ctx context.Context, client dynamodb.Client, url url.URL) (string, error) {
	Slug := uuid.New().String()
	Url := url.String()

	item, err := attributevalue.MarshalMap(ShortenedUrlDocument{
		Slug,
		Url,
	})

	if err != nil {
		return "", fmt.Errorf("Can't add url to db: %w", err)
	}

	tableName := os.Getenv("URLS_TABLE_NAME")

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	})

	if err != nil {
		return "", fmt.Errorf("Put in table '%s' failed: %w", tableName, err)
	}

	return Slug, nil
}
