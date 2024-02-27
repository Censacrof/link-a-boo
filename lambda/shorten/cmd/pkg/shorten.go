package shorten

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
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
