package util

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type LambdaResponse events.APIGatewayProxyResponse
type LambdaRequest events.APIGatewayProxyRequest

func MakeRespose(statusCode int, body interface{}) LambdaResponse {
	json, _ := json.Marshal(body)

	return LambdaResponse{
		StatusCode: statusCode,
		Body:       string(json),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}
}
