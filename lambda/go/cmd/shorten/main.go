package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"net/url"

	"github.com/Censacrof/link-a-boo/lambda/shorten/pkg/db"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

	shortenedUrl := db.NewShortenedUrl(*targetUrl)
	ddbClient := dynamodb.NewFromConfig(cfg)

	err = shortenedUrl.Put(ctx, *ddbClient)

	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Internal server error: %v", err),
		}, nil
	}

	resp := ShortenResponse{
		ShortenedUrl: shortenedUrl.Slug,
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

func main() {
	lambda.Start(HandleShortenRequest)
}
