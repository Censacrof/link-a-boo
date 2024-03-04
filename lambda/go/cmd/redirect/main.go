package main

import (
	"context"
	"fmt"

	"github.com/Censacrof/link-a-boo/lambda/go/pkg/db/shortened_url"
	"github.com/Censacrof/link-a-boo/lambda/go/pkg/response"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRedirectRequest(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	slug, ok := event.PathParameters["slug"]
	if !ok {
		return response.NewErrorResponse("Slug missing in request path").ToApiGatewayProxyResponse(400)
	}

	shortenedUrl, err := shortened_url.Get(ctx, slug)
	if err != nil {
		return response.NewErrorResponse(fmt.Sprintf("Internal server error: %v", err)).ToApiGatewayProxyResponse(500)
	}

	if shortenedUrl == nil {
		return response.NewErrorResponse(fmt.Sprintf("Can't find URL corresponding to slug: %s", slug)).ToApiGatewayProxyResponse(404)
	}

	resp, err := response.NewOkResponse(fmt.Sprintf("Slug: %s", slug)).ToApiGatewayProxyResponse(301)
	if err != nil {
		return resp, err
	}

	resp.Headers["Location"] = shortenedUrl.Url
	resp.Headers["Cache-Control"] = "max-age=180, private"

	return resp, err
}

func main() {
	lambda.Start(HandleRedirectRequest)
}
