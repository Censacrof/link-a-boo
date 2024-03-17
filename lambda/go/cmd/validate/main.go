package main

import (
	"context"

	"github.com/Censacrof/link-a-boo/lambda/go/pkg/response"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleValidateRequest(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return response.NewOkResponse("valid").ToApiGatewayProxyResponse(200)
}

func main() {
	lambda.Start(HandleValidateRequest)
}
