package main

import (
	"context"
	"encoding/json"
	"errors"
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

	slug := validateRequest.Slug
	if slug == "" {
		slug = "aValidSlug"
	}

	_, err = shortened_url.New(validateRequest.Url, slug)
	if err != nil {
		if errors.Is(err, shortened_url.ErrInvalidUrl) {
			return response.NewErrorResponse("invalid url").ToApiGatewayProxyResponse(400)
		}

		if errors.Is(err, shortened_url.ErrInvalidSlug) {
			return response.NewErrorResponse("invalid slug").ToApiGatewayProxyResponse(400)
		}

		return response.NewErrorResponse("invalid").ToApiGatewayProxyResponse(400)
	}

	if validateRequest.Slug != "" {
		slugAvailable, err := shortened_url.IsSlugAvailable(ctx, validateRequest.Slug)
		if err != nil {
			return response.NewErrorResponse("internal server error").ToApiGatewayProxyResponse(500)
		}

		if !slugAvailable {
			return response.NewErrorResponse("slug already in use").ToApiGatewayProxyResponse(409)
		}
	}

	return response.NewOkResponse("valid").ToApiGatewayProxyResponse(200)
}

func main() {
	lambda.Start(HandleValidateRequest)
}
