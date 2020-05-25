package main

import (
	"context"
  "encoding/json"
  "fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/igoratron/blend/pkgs/store"
)

type httpResponse events.APIGatewayProxyResponse
type httpRequest events.APIGatewayProxyRequest

func GetIngredients(ctx context.Context, event httpRequest) (httpResponse, error) {
  ingredientName := event.QueryStringParameters["q"]
  ingredients, err := store.SearchIngredients(&ingredientName)

  if ingredients == nil {
    ingredients = []store.Ingredient{}
  }

  if err != nil {
    fmt.Println(err)
    return makeRespose(500, err), nil
  }

  return makeRespose(200, ingredients), nil
}

func makeRespose(statusCode int, body interface{}) httpResponse {
  json, _ := json.Marshal(body)

  return httpResponse {
    StatusCode: 200,
    Body: string(json),
    Headers: map[string]string{
      "Access-Control-Allow-Origin": "*",
    },
  }
}

func main() {
	lambda.Start(GetIngredients)
}
