package response

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
)

type ProblemDetails struct {
	Type     url.URL `json:"type"`
	Status   int     `json:"status,omitempty"`
	Title    string  `json:"title,omitempty"`
	Detail   string  `json:"detail,omitempty"`
	Instance url.URL `json:"instance,omitempty"`
}

func (self *ProblemDetails) ToApiGatewayProxyResponse() (*events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(self)
	if err != nil {
		return &events.APIGatewayProxyResponse{}, fmt.Errorf("Can't convert to API gateway response: %w", err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: self.Status,
		Body:       string(body),
	}, nil
}
