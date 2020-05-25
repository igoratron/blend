package main

import (
	"context"
  "encoding/json"
  "fmt"
  "strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/igoratron/blend/pkgs/store"
)

type httpResponse events.APIGatewayProxyResponse
type httpRequest events.APIGatewayProxyRequest

func GetRecipes(ctx context.Context, event httpRequest) (httpResponse, error) {
  ingredientIds := strings.Split(event.QueryStringParameters["ingredientIds"], ",")
  recipes, err := store.GetRecipes(&ingredientIds)

  if err != nil {
    fmt.Println(err)
    return makeRespose(500, err), nil
  }

  return makeRespose(200, recipes), nil
}

func makeRespose(statusCode int, body interface{}) httpResponse {
  json, _ := json.Marshal(body)

  return httpResponse {
    StatusCode: statusCode,
    Body: string(json),
    Headers: map[string]string{
      "Access-Control-Allow-Origin": "*",
    },
  }
}

func main() {
	lambda.Start(GetRecipes)
}
