import React, { useCallback, useEffect, useState } from 'react';
import './App.css';

function searchIngredients(ingredientName) {
  const apiUrl = "https://mfuqctb1me.execute-api.eu-west-1.amazonaws.com/dev/ingredients";
  return fetch(`apiUrl?q=${ingredientName}`).then((response) => response.json());
}

function IngredientSearch() {
  const [ingredients, setIngredients] = useState([]);
  const onInputChange = useCallback(event => {
    searchIngredients(event.currentTarget.value)
      .then(setIngredients)
      .catch(() => setIngredients([]));
  }, []);

  return (
    <>
      <input placeholder="Add ingredient" onChange={onInputChange} />
      <ul className="list">
        {ingredients.map(ingredient => (
          <li key={ingredient.id}>
            <label>
              <input type="checkbox" value="123" />
              {ingredient.name}
            </label>
          </li>
        ))}
      </ul>
    </>
  );
}

function App() {
  return (
    <div className="two-panes">
      <div className="two-panes_left">
        <h2 className="title">Ingredients you have</h2>
        <ul className="list ingredients">
          <li>Onion</li>
        </ul>
        <IngredientSearch />
      </div>
      <div className="two-panes_right">
        <h2 className="title">Available recipes</h2>
        <ul className="list">
          <li>Onion</li>
        </ul>
      </div>
    </div>
  );
}

export default App;
