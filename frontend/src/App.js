import React, { useReducer, useEffect } from "react";
import "./App.css";

import * as actions from "./actions";
import IngredientSearch from "./ingredienSearch";

export const StateContext = React.createContext();

const initialState = {
  pantry: [],
  recipes: [],
  ingredientIdQuery: ""
};

const reducer = (state, action) => {
  let pantry;
  switch (action.type) {
    case actions.INGREDIENT_ADDED:
      pantry = state.pantry.concat(action.payload);
      return {
        ...state,
        pantry
      };
    case actions.INGREDIENT_REMOVED:
      const remove = action.payload;
      pantry = state.pantry.filter(i => remove.id !== i.id);
      return {
        ...state,
        pantry
      };
    case actions.RECIPES_REQUESTED:
      const ingredientIds = state.pantry.map(i => i.id).join(",");
      return {
        ...state,
        ingredientIdQuery: ingredientIds
      };
    case actions.RECIPES_FETCHED:
      return {
        ...state,
        recipes: action.payload
      };
    default:
      console.log("Unknown action:", action);
      return state;
  }
};

function App() {
  const [state, dispatch] = useReducer(reducer, initialState);

  const { pantry, recipes, ingredientIdQuery } = state;

  useEffect(() => {
    if (!ingredientIdQuery) {
      return;
    }
    const apiUrl =
      "https://mfuqctb1me.execute-api.eu-west-1.amazonaws.com/dev/recipes";
    fetch(`${apiUrl}?ingredientIds=${ingredientIdQuery}`)
      .then(response => response.json())
      .then(recipes => dispatch(actions.recipesFetched(recipes)))
      .catch(() => dispatch(actions.recipesFetched([])));
  }, [ingredientIdQuery]);

  return (
    <StateContext.Provider value={{ state, dispatch }}>
      <div className="two-panes">
        <div className="two-panes_left">
          <h2 className="title">Ingredients you have</h2>
          <ul className="list ingredients">
            {pantry.map(({ id, name }) => (
              <li key={id}>{name}</li>
            ))}
          </ul>
          <IngredientSearch />
        </div>
        <div className="two-panes_right">
          <h2 className="title">Available recipes</h2>
          <ul className="list">
            {recipes.map((recipe, i) => (
              <li key={i}>
                <a href={recipe.url} target="_blank" rel="noopener noreferrer">
                  {recipe.name}
                </a>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </StateContext.Provider>
  );
}

export default App;
