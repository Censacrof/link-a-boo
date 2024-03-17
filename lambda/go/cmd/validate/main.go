package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Censacrof/link-a-boo/lambda/go/pkg/db/shortened_url"
	"github.com/Censacrof/link-a-boo/lambda/go/pkg/response"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ValidateRequest struct {
	Url  string `json:"Url"`
	Slug string `json:"slug"`
}

func HandleValidateRequest(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var validateRequest ValidateRequest
	err := json.Unmarshal([]byte(event.Body), &validateRequest)
	if err != nil {
		return response.NewErrorResponse(fmt.Sprintf("Invalid request: %v", err)).ToApiGatewayProxyResponse(400)
	}

	_, err = shortened_url.New(validateRequest.Url, validateRequest.Slug)
	if err != nil {
		return response.NewErrorResponse("invalid").ToApiGatewayProxyResponse(400)
	}

	return response.NewOkResponse("valid").ToApiGatewayProxyResponse(200)
}

func main() {
	lambda.Start(HandleValidateRequest)
}
