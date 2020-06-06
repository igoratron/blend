package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/igoratron/blend/pkgs/store"
	u "github.com/igoratron/blend/pkgs/util"
)

func GetShoppingList(ctx context.Context, event u.LambdaRequest) (u.LambdaResponse, error) {
	ingredientIds := strings.Split(event.QueryStringParameters["recipeIds"], ",")

	recipes, err := store.GetShoppingList(ingredientIds)

	if err != nil {
		fmt.Println(err)
		return u.MakeRespose(500, err), nil
	}

	return u.MakeRespose(200, recipes), nil
}

func main() {
	lambda.Start(GetShoppingList)
}
