package store

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func toInterfaceArray(array []string) []interface{} {
	result := make([]interface{}, len(array))

	for i, e := range array {
		result[i] = e
	}

	return result
}

func GetRecipes(ingredientIds []string) ([]Recipe, error) {
	placeholderString := makePlaceholderString(len(ingredientIds), 1)
	sqlStatement := fmt.Sprintf(`
    SELECT recipe_id, count(*) c
    FROM ingredients_recipes
    WHERE ingredient_id IN %s
    GROUP BY recipe_id
    ORDER BY c DESC
    LIMIT 20
  `, placeholderString)

	fmt.Println("Querying:", sqlStatement, ingredientIds)

	var recipes []Recipe
	rows, err := db.Query(sqlStatement, toInterfaceArray(ingredientIds)...)
	if err != nil {
		return recipes, err
	}

	var recipeIds []string
	for rows.Next() {
		var recipeId string
		var count int
		rows.Scan(&recipeId, &count)
		recipeIds = append(recipeIds, recipeId)
	}

	fmt.Printf("Found %d recipies in DB:\n%#v\n", len(recipeIds), recipeIds)

	return getRecipesFromDDB(recipeIds)
}

func getRecipesFromDDB(recipeIds []string) ([]Recipe, error) {
	svc := dynamodb.New(session.New())

	keys := make([]map[string]*dynamodb.AttributeValue, len(recipeIds))
	for i, recipeId := range recipeIds {
		keys[i] = map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(recipeId),
			},
		}
	}

	requestItems := map[string]*dynamodb.KeysAndAttributes{
		"hellofresh-recipes": {
			Keys: keys,
		},
	}

	batchQuery := dynamodb.BatchGetItemInput{
		RequestItems: requestItems,
	}

	output, err := svc.BatchGetItem(&batchQuery)

	if err != nil {
		return []Recipe{}, err
	}

	recipes := output.Responses["hellofresh-recipes"]

	result := make([]Recipe, len(recipes))

	for i, recipe := range recipes {
		result[i] = Recipe{
			Name: *recipe["name"].S,
			Url:  *recipe["websiteUrl"].S,
		}
	}

	return result, nil
}
