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
	Url  string `json:"Url" validate:"required,min=5"`
	Slug string `json:"slug"`
}

type ShortenResponse struct {
	Slug string `json:"slug"`
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

	slug := shortenRequest.Slug
	if shortenRequest.Slug == "" {
		slug = uuid.New().String()
	}

	shortenedUrl, err := shortened_url.New(shortenRequest.Url, slug)
	if err != nil {
		return response.NewErrorResponse(fmt.Sprintf("Invalid request: %v", err)).ToApiGatewayProxyResponse(400)
	}

	err = shortened_url.Put(ctx, shortenedUrl)
	if err != nil {
		if errors.Is(err, shortened_url.ErrSlugAlreadyExists) {
			return response.NewErrorResponse("Slug already exists").ToApiGatewayProxyResponse(409)
		}

		return response.NewErrorResponse(fmt.Sprintf("Internal server error: %v", err)).ToApiGatewayProxyResponse(500)
	}

	return response.NewOkResponse(ShortenResponse{
		Slug: shortenedUrl.Slug,
	}).ToApiGatewayProxyResponse(201)
}

func main() {
	lambda.Start(HandleShortenRequest)
}
