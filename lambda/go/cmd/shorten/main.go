package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Censacrof/link-a-boo/lambda/go/pkg/db/shortened_url"
	"github.com/Censacrof/link-a-boo/lambda/go/pkg/response"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"

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

	slug := uuid.New()
	shortenedUrl, err := shortened_url.New(shortenRequest.TargetUrl, slug.String())
	if err != nil {
		return response.NewErrorResponse(fmt.Sprintf("Invalid request: %v", err)).ToApiGatewayProxyResponse(400)
	}

	err = shortened_url.Put(ctx, shortenedUrl)

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
