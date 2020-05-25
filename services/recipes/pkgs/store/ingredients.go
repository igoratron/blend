package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

var db *sql.DB

func init() {
	if db == nil {
		db = connect()
	}
}

func AddIngredients(recipeId *string, ingredientNames *[]string) error {
	ingredientIds := make([]string, len(*ingredientNames))
	for index, name := range *ingredientNames {
		ingredient := Ingredient{
			Id:   uuid.New().String(),
			Name: name,
		}

		fmt.Printf("Adding %s\n", name)
		err := insertInto(db, "ingredients", columns("id", "name"), &[]Record{&ingredient})
		if err != nil {
			ingredient, err = GetIngredient(name)
			fmt.Printf("Already exists with id %s\n", ingredient.Id)
		}
		if err != nil {
			return err
		}
		ingredientIds[index] = ingredient.Id
	}

	ingredientRecipes := make([]Record, len(*ingredientNames))

	for index, ingredientId := range ingredientIds {
		ingredientRecipes[index] = IngredientRecipe{
			IngredientId: ingredientId,
			RecipeId:     *recipeId,
		}
	}

	err := insertInto(db, "ingredients_recipes", columns("ingredient_id", "recipe_id"), &ingredientRecipes)

	if err != nil {
		return err
	}

	fmt.Println("All done.")

	return nil
}

func GetIngredient(name string) (Ingredient, error) {
	var ingredient Ingredient

	row := selectFrom(db, columns("id", "name"), "ingredients", "name = ?", name)
	err := row.Scan(&ingredient.Id, &ingredient.Name)

	if err != nil {
		return ingredient, err
	}

	return ingredient, nil
}

func SearchIngredients(ingredientName *string) ([]Ingredient, error) {
	var ingredients []Ingredient
	fmt.Printf("Searching for %s\n", *ingredientName)
	wildcard := *ingredientName + "*"
	results, err := db.Query(`SELECT id, name FROM ingredients WHERE MATCH(name) AGAINST (? IN BOOLEAN MODE)`, wildcard)

	if err != nil {
		return ingredients, err
	}

	for results.Next() {
		var ingredient Ingredient
		err = results.Scan(&ingredient.Id, &ingredient.Name)

		if err != nil {
			return ingredients, err
		}

		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}
