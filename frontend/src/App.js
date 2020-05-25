import React, { useCallback, useState } from "react";
import "./App.css";

function searchIngredients(ingredientName) {
  const apiUrl =
    "https://mfuqctb1me.execute-api.eu-west-1.amazonaws.com/dev/ingredients";
  return fetch(`${apiUrl}?q=${ingredientName}`).then(response =>
    response.json()
  );
}

function IngredientSearch({ pantry, addToPantry, removeFromPantry }) {
  const [ingredients, setIngredients] = useState([]);
  const onInputChange = useCallback(event => {
    searchIngredients(event.currentTarget.value)
      .then(setIngredients)
      .catch(() => setIngredients([]));
  }, []);
  const onListItemClick = useCallback(
    ingredient => event => {
      const isChecked = event.target.checked;

      if (isChecked) {
        addToPantry(ingredient);
      } else {
        removeFromPantry(ingredient);
      }
    },
    [addToPantry, removeFromPantry]
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

function App() {
  const [recipes, setRecipes] = useState([]);
  const [pantry, setPantry] = useState([]);
  const searchRecipes = useCallback(
    pantry => {
      const ingredientIds = pantry.map(i => i.id).join(",");
      const apiUrl =
        "https://mfuqctb1me.execute-api.eu-west-1.amazonaws.com/dev/recipes";
      return fetch(`${apiUrl}?ingredientIds=${ingredientIds}`)
        .then(response => response.json())
        .then(setRecipes)
        .catch(() => setRecipes([]));
    },
    [pantry]
  );
  const addToPantry = ingredient => {
    const newPantry = pantry.concat(ingredient);
    setPantry(newPantry);
    searchRecipes(newPantry);
  };
  const removeFromPantry = ingredient => {
    setPantry(pantry.filter(i => ingredient.id !== i.id));
  };

  return (
    <div className="two-panes">
      <div className="two-panes_left">
        <h2 className="title">Ingredients you have</h2>
        <ul className="list ingredients">
          {pantry.map(({ id, name }) => (
            <li key={id}>{name}</li>
          ))}
        </ul>
        <IngredientSearch
          pantry={pantry}
          addToPantry={addToPantry}
          removeFromPantry={removeFromPantry}
        />
      </div>
      <div className="two-panes_right">
        <h2 className="title">Available recipes</h2>
        <ul className="list">
          {recipes.map(recipe => (
            <li>
              <a href={recipe.url} target="_blank" rel="noopener noreferrer">{recipe.name}</a>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}

export default App;
