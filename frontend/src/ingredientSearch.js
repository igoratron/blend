import React, { useContext, useCallback, useState } from "react";

import { StateContext } from "./App";
import * as actions from "./actions";

function searchIngredients(ingredientName) {
  const apiUrl =
    "https://mfuqctb1me.execute-api.eu-west-1.amazonaws.com/dev/ingredients";
  return fetch(`${apiUrl}?q=${ingredientName}`).then(response =>
    response.json()
  );
}

export default function IngredientSearch() {
  const [ingredients, setIngredients] = useState([]);
  const { state, dispatch } = useContext(StateContext);
  const { pantry } = state;

  const onInputChange = useCallback(event => {
    searchIngredients(event.currentTarget.value)
      .then(setIngredients)
      .catch(() => setIngredients([]));
  }, []);

  const onListItemClick = useCallback(
    ingredient => event => {
      const isChecked = event.target.checked;

      if (isChecked) {
        dispatch(actions.addIngredient(ingredient));
      } else {
        dispatch(actions.removeIngredient(ingredient));
      }
      dispatch(actions.suggestRecipes());
    },
    [dispatch]
  );

  return (
    <>
      <input placeholder="Add ingredient" onChange={onInputChange} />
      <ul className="list">
        {ingredients.map(ingredient => (
          <li key={ingredient.id}>
            <label className="list_item">
              <input
                type="checkbox"
                value={ingredient.id}
                onChange={onListItemClick(ingredient)}
                checked={pantry.some(i => i.id === ingredient.id)}
              />
              {ingredient.name}
            </label>
          </li>
        ))}
      </ul>
    </>
  );
}
