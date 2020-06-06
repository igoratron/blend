export const RECIPES_FETCHED = "RECIPES_FETCHED";
export const RECIPES_REQUESTED = "RECIPES_REQUESTED";
export const ADDED_TO_PLAN = "ADDED_TO_PLAN";
export const REMOVED_FROM_PLAN = "REMOVED_FROM_PLAN";

export const suggestRecipes = ingredientIds => ({
  type: RECIPES_REQUESTED,
  payload: ingredientIds
});

export const recipesFetched = recipes => ({
  type: RECIPES_FETCHED,
  payload: recipes
});

export const recipeAddedToPlan = recipe => ({
  type: ADDED_TO_PLAN,
  payload: recipe
});

export const recipeRemovedFromPlan = recipe => ({
  type: REMOVED_FROM_PLAN,
  payload: recipe
});
