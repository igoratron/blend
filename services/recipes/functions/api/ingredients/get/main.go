package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/igoratron/blend/pkgs/store"
	u "github.com/igoratron/blend/pkgs/util"
)

func GetIngredients(ctx context.Context, event u.LambdaRequest) (u.LambdaResponse, error) {
	ingredientName := event.QueryStringParameters["q"]
	ingredients, err := store.SearchIngredients(ingredientName)

	if err != nil {
		fmt.Println(err)
		return u.MakeRespose(500, err), nil
	}

	return u.MakeRespose(200, ingredients), nil
}

func main() {
	lambda.Start(GetIngredients)
}
