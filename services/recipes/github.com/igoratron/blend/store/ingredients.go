package store

import (
  "fmt"
	"database/sql"

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
    err := insertInto(db, "ingredients", columns("id", "name"), &[]Record{ &ingredient })
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
			IngredientId:   ingredientId,
			RecipeId: *recipeId,
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

func GetRecommendedRecipes() (*[]string, error) {
	sqlStatement := `
    SELECT recipe_id, count(*) c
    FROM pantry p JOIN ingredients_recipes ir
    ON p.ingredient_id = ir.ingredient_id
    GROUP BY recipe_id
    ORDER BY c DESC
    LIMIT 3
  `
  recipeIds := make([]string, 3)
  rows, err := db.Query(sqlStatement)
  if err != nil {
    return &recipeIds, err
  }

  for i := 0; rows.Next(); i += 1 {
    var count int
    rows.Scan(&recipeIds[i], &count)
  }

  return &recipeIds, nil
}
