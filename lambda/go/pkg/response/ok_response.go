package response

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type OkResponse[T any] struct {
	Ok T `json:"ok"`
}

func NewOkResponse[T any](data T) *OkResponse[T] {
	return &OkResponse[T]{
		Ok: data,
	}
}

func (self *OkResponse[T]) ToApiGatewayProxyResponse(statusCode int) (*events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(self)
	if err != nil {
		return &events.APIGatewayProxyResponse{}, fmt.Errorf("Can't convert to API gateway response: %w", err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
	}, nil
}
