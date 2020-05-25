package main

import (
	"context"
  "fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

  "github.com/igoratron/blend/store"
)

type Event events.DynamoDBEvent

func Handler(ctx context.Context, e Event) {
	for _, record := range e.Records {
		if record.EventName != "INSERT" {
			fmt.Printf("Event type is %s, skipping...", record.EventName)
			return
		}

		newRecipe := record.Change.NewImage
		recipeId := newRecipe["id"].String()
		recipeIngredients := extractIngredientNames(newRecipe["ingredients"].List())

		fmt.Printf("Adding recipe %s. Found %d ingredients\n", recipeId, len(*recipeIngredients))

    err := store.AddIngredients(&recipeId, recipeIngredients)

    if err != nil {
      fmt.Println(err)
      continue
    }

    fmt.Println("Recipe added")
	}
}

func extractIngredientNames(ingredients []events.DynamoDBAttributeValue) *[]string {
	result := make([]string, len(ingredients))

	for index, ingredient := range ingredients {
		result[index] = ingredient.Map()["name"].String()
	}

	return &result
}

func main() {
	lambda.Start(Handler)
}
