package main

import (
	"context"
	"fmt"

	"github.com/Censacrof/link-a-boo/lambda/go/pkg/db/shortened_url"
	"github.com/Censacrof/link-a-boo/lambda/go/pkg/response"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type redirectResult struct {
	Location string `json:"location"`
}

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

	resp, err := response.NewOkResponse(redirectResult{Location: shortenedUrl.Url.String()}).ToApiGatewayProxyResponse(301)
	if err != nil {
		return resp, err
	}

	resp.Headers = make(map[string]string)
	resp.Headers["Location"] = shortenedUrl.Url.String()
	resp.Headers["Cache-Control"] = "max-age=180, private"

	return resp, nil
}

func main() {
	lambda.Start(HandleRedirectRequest)
}
