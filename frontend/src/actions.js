export const INGREDIENT_ADDED = "INGREDIENT_ADDED";
export const INGREDIENT_QUERIED = "INGREDIENT_QUERIED";
export const INGREDIENT_REMOVED = "INGREDIENT_REMOVED";
export const RECIPES_FETCHED = "RECIPES_FETCHED";
export const RECIPES_REQUESTED = "RECIPES_REQUESTED";

export const addIngredient = ingredient => ({
  type: INGREDIENT_ADDED,
  payload: ingredient
});

export const removeIngredient = ingredient => ({
  type: INGREDIENT_REMOVED,
  payload: ingredient
});

export const ingredientsFetched = ingredients => ({
  type: INGREDIENT_QUERIED,
  payload: ingredients
});

export const suggestRecipes = () => ({
  type: RECIPES_REQUESTED,
});

export const recipesFetched = recipes => ({
  type: RECIPES_FETCHED,
  payload: recipes
});
