package response

import "github.com/aws/aws-lambda-go/events"

type Response interface {
	ToApiGatewayProxyResponse(statusCode int) (*events.APIGatewayProxyResponse, error)
}
