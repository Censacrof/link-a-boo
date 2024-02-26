package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	message := fmt.Sprintf("Hello World! _HANDLER=%s; %+v", os.Getenv("_HANDLER"), *event)

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       message,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
