package response

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorResponseError struct {
	Reason string `json:"reason"`
}

type ErrorResponse struct {
	Error ErrorResponseError `json:"error"`
}

func NewErrorResponse(reason string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorResponseError{
			Reason: reason,
		},
	}
}

func (self *ErrorResponse) ToApiGatewayProxyResponse(statusCode int) (*events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(self)
	if err != nil {
		return &events.APIGatewayProxyResponse{}, fmt.Errorf("Can't convert to API gateway response: %w", err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
	}, nil
}
