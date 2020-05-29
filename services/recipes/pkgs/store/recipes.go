package store

import (
	"fmt"
	"sort"

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
	var ingredientCount []int
	for rows.Next() {
		var recipeId string
		var count int
		rows.Scan(&recipeId, &count)
		recipeIds = append(recipeIds, recipeId)
		ingredientCount = append(ingredientCount, count)
	}

	fmt.Printf("Found %d recipies in DB:\n%#v\n", len(recipeIds), recipeIds)

	recipeDetails, err := getRecipesFromDDB(recipeIds)

	if err != nil {
		return recipes, err
	}

	recipes = make([]Recipe, len(recipeIds))
	for index, recipeId := range recipeIds {
		rd := recipeDetails[recipeId]
		rd.Ingredients.Matching = ingredientCount[index]
		recipes[index] = rd
	}

	sort.Slice(recipes, func(i, j int) bool {
		iMatch := float32(recipes[i].Ingredients.Matching) / float32(recipes[i].Ingredients.Total)
		jMatch := float32(recipes[j].Ingredients.Matching) / float32(recipes[j].Ingredients.Total)

		return iMatch > jMatch
	})

	return recipes, nil
}

func getRecipesFromDDB(recipeIds []string) (map[string]Recipe, error) {
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
		return map[string]Recipe{}, err
	}

	recipes := output.Responses["hellofresh-recipes"]

	result := make(map[string]Recipe)

	for _, recipe := range recipes {
		result[*recipe["id"].S] = Recipe{
			Name: *recipe["name"].S,
			Url:  *recipe["websiteUrl"].S,
			Ingredients: IngredientMatch{
				Total: len(recipe["ingredients"].L),
			},
		}
	}

	return result, nil
}
