export const RECIPES_FETCHED = "RECIPES_FETCHED";
export const RECIPES_REQUESTED = "RECIPES_REQUESTED";

export const suggestRecipes = (ingredientIds) => ({
  type: RECIPES_REQUESTED,
  payload: ingredientIds
});

export const recipesFetched = recipes => ({
  type: RECIPES_FETCHED,
  payload: recipes
});
