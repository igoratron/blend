package store

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/igoratron/blend/pkgs/hellofresh"
)

func toInterfaceArray(array []string) []interface{} {
	result := make([]interface{}, len(array))

	for i, e := range array {
		result[i] = e
	}

	return result
}

func getIngredientNames(ingredients []hellofresh.Ingredient) []string {
	names := make([]string, len(ingredients))

	for index, ingredient := range ingredients {
		names[index] = ingredient.Name
	}

	return names
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

	var recipeIds []hellofresh.Id
	var ingredientCount []int
	for rows.Next() {
		var recipeId hellofresh.Id
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
		hfRecipe := recipeDetails[recipeId]
		recipes[index] = Recipe{
			Id:   string(hfRecipe.Id),
			Name: hfRecipe.Name,
			Url:  hfRecipe.WebsiteUrl,
			Ingredients: Ingredients{
				Matching: ingredientCount[index],
				Total:    len(hfRecipe.Ingredients),
				List:     getIngredientNames(hfRecipe.Ingredients),
			},
			ImagePath: hfRecipe.ImagePath,
		}
	}

	sort.Slice(recipes, func(i, j int) bool {
		iMatch := float32(recipes[i].Ingredients.Matching) / float32(recipes[i].Ingredients.Total)
		jMatch := float32(recipes[j].Ingredients.Matching) / float32(recipes[j].Ingredients.Total)

		return iMatch > jMatch
	})

	return recipes, nil
}

func getRecipesFromDDB(recipeIds []hellofresh.Id) (map[hellofresh.Id]hellofresh.Recipe, error) {
	svc := dynamodb.New(session.New())
	result := make(map[hellofresh.Id]hellofresh.Recipe)
	var err error
	fmt.Println("Getting recipes from dynamo")

	keys := make([]map[string]*dynamodb.AttributeValue, len(recipeIds))
	for i, recipeId := range recipeIds {
		keys[i] = map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(string(recipeId)),
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
		return result, err
	}

	hellofreshRecipes := output.Responses["hellofresh-recipes"]

	recipes := make([]hellofresh.Recipe, len(hellofreshRecipes))
	for index, hellofreshRecipe := range hellofreshRecipes {
		recipes[index] = makeHelloFreshRecipe(hellofreshRecipe)
	}

	for _, recipe := range recipes {
		result[hellofresh.Id(recipe.Id)] = recipe
	}

	return result, nil
}

func makeHelloFreshRecipe(entry map[string]*dynamodb.AttributeValue) hellofresh.Recipe {
	recipe := hellofresh.Recipe{
		Id:         *entry["id"].S,
		Name:       *entry["name"].S,
		WebsiteUrl: *entry["websiteUrl"].S,
		ImagePath:  *entry["imagePath"].S,
	}

	errs := []error{}

	ingredients := []hellofresh.Ingredient{}
	err := dynamodbattribute.UnmarshalList(entry["ingredients"].L, &ingredients)
	if err != nil {
		errs = append(errs, err)
	}

	yields := []hellofresh.Yield{}
	err = dynamodbattribute.UnmarshalList(entry["yields"].L, &yields)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		fmt.Println(errs)
	} else {
		recipe.Ingredients = ingredients
		recipe.Yields = yields
	}

	return recipe
}
