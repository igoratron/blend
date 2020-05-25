package main

import (
	"context"
  "encoding/json"
  "fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/igoratron/blend/store"
)

type RecipeDetail struct {
  Name string `json:"name"`
  Url string `json:"url"`
}

type ShoppingListItem struct {
  Name string `json:"name"`
  Amount string `json:"amount"`
}

type ShoppingList struct {
  Recipes []RecipeDetail `json:"recipes"`
  ShoppingList []ShoppingListItem `json:"shoppingList"`
}

type httpResponse events.APIGatewayProxyResponse


func GetShoppingList(ctx context.Context) (httpResponse, error) {
  recipeIds, err := store.GetRecommendedRecipes()
  if err != nil {
    return makeRespose(500, err.Error()), nil
  }

  shoppingList := getRecipies(recipeIds)

  return makeRespose(200, shoppingList), nil
}

func getRecipies(recipeIds *[]string) ShoppingList {
  var shoppingList ShoppingList
  svc := dynamodb.New(session.New())

  keys := make([]map[string]*dynamodb.AttributeValue, len(*recipeIds))
  for i, recipeId := range *recipeIds {
    fmt.Printf("key %d: %s\n", i, recipeId)
    keys[i] = map[string]*dynamodb.AttributeValue{
      "id": {
        S: aws.String(recipeId),
      },
    }
  }

  requestItems := map[string]*dynamodb.KeysAndAttributes {
    "hellofresh-recipes": {
      Keys: keys,
    },
  }

  batchQuery := dynamodb.BatchGetItemInput{
    RequestItems: requestItems,
  }

  output, err := svc.BatchGetItem(&batchQuery)

  if err != nil {
    fmt.Println("Error:")
    fmt.Println(err)
    return shoppingList
  }

  recipes := output.Responses["hellofresh-recipes"]

  shoppingList.Recipes = make([]RecipeDetail, len(recipes))

  for i, recipe := range recipes {
    shoppingList.Recipes[i] = RecipeDetail{
      Name: *recipe["name"].S,
      Url: *recipe["websiteUrl"].S,
    }
  }

  return shoppingList
}

func makeRespose(statusCode int, body interface{}) httpResponse {
  json, _ := json.Marshal(body)

  return httpResponse {
    StatusCode: 200,
    Body: string(json),
  }
}

func main() {
	lambda.Start(GetShoppingList)
}
