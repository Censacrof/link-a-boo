package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ShortenRequest struct {
	TargetUrl string `json:"targetUrl"`
}

type ShortenResponse struct {
	ShortenedUrl string `json:"shortenedUrl"`
}

func HandleShortenRequest(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var shortenRequest ShortenRequest
	err := json.Unmarshal([]byte(event.Body), &shortenRequest)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request",
		}, nil
	}

	resp := ShortenResponse{
		ShortenedUrl: "shortened-" + shortenRequest.TargetUrl,
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

func getHandler() (interface{}, error) {
	handlerName := os.Getenv("_HANDLER")

	switch handlerName {
	case "shorten":
		return HandleShortenRequest, nil

	default:
		return nil, errors.New(fmt.Sprintf("Invalid handler %s", handlerName))
	}
}

func main() {
	handler, err := getHandler()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	lambda.Start(handler)
}
