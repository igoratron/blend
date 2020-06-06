package store

import (
	"github.com/igoratron/blend/pkgs/hellofresh"
)

func makeHelloFreshArray(hfids []string) []hellofresh.Id {
	result := make([]hellofresh.Id, len(hfids))

	for index, hfid := range hfids {
		result[index] = hellofresh.Id(hfid)
	}

	return result
}

func GetShoppingList(recipeIds []string) (ShoppingList, error) {
	var result ShoppingList
	hfids := makeHelloFreshArray(recipeIds)
	recipes, err := getRecipesFromDDB(hfids)
	if err != nil {
		return result, err
	}

	ingredients := make(map[hellofresh.Id]hellofresh.Ingredient)
	shoppingList := make(map[hellofresh.Id]ShoppingListItem)

	for _, recipe := range recipes {
		for _, ingredient := range recipe.Ingredients {
			ingredients[hellofresh.Id(ingredient.Id)] = ingredient
		}

		ingredientAmounts := recipe.Yields[len(recipe.Yields)-1].Ingredients

		for _, ingredientAmount := range ingredientAmounts {
			iaId := hellofresh.Id(ingredientAmount.Id)
			if slItem, ok := shoppingList[iaId]; ok {
				slItem.Amount += ingredientAmount.Amount
			} else {
				shoppingList[iaId] = ShoppingListItem{
					Name:   ingredients[iaId].Name,
					Amount: ingredientAmount.Amount,
					Unit:   ingredientAmount.Unit,
				}
			}
		}
	}

	for _, slItem := range shoppingList {
		result.Ingredients = append(result.Ingredients, slItem)
	}

	return result, nil
}
