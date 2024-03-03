package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/Censacrof/link-a-boo/lambda/go/pkg/db"
	"github.com/Censacrof/link-a-boo/lambda/go/pkg/response"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-lambda-go/events"

	"github.com/go-playground/validator/v10"
)

type ShortenRequest struct {
	TargetUrl string `json:"targetUrl" validate:"required,min=5"`
}

type ShortenResponse struct {
	ShortenedUrl string `json:"shortenedUrl"`
}

func HandleShortenRequest(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var shortenRequest ShortenRequest
	unmarshalErr := json.Unmarshal([]byte(event.Body), &shortenRequest)

	validate := validator.New(validator.WithRequiredStructEnabled())
	validationErr := validate.Struct(shortenRequest)

	err := errors.Join(unmarshalErr, validationErr)
	if err != nil {
		return response.NewErrorResponse(fmt.Sprintf("Invalid request: %v", err)).ToApiGatewayProxyResponse(400)
	}

	targetUrl, err := url.Parse(shortenRequest.TargetUrl)
	if err != nil {
		return response.NewErrorResponse("Invalid targetUrl").ToApiGatewayProxyResponse(400)
	}

	shortenedUrl := db.NewShortenedUrl(*targetUrl)
	err = db.GetShortenedUrlTable().Put(ctx, shortenedUrl)

	if err != nil {
		return response.NewErrorResponse(fmt.Sprintf("Internal server error: %v", err)).ToApiGatewayProxyResponse(500)
	}

	return response.NewOkResponse(ShortenResponse{
		ShortenedUrl: shortenedUrl.Slug,
	}).ToApiGatewayProxyResponse(200)
}

func main() {
	lambda.Start(HandleShortenRequest)
}
