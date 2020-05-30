import React, { useContext, useCallback, useState } from "react";
import { Select } from 'antd';

import { StateContext } from "./App";
import * as actions from "./actions";

const { Option } = Select;

function searchIngredients(ingredientName, callback) {
  const apiUrl =
    "https://mfuqctb1me.execute-api.eu-west-1.amazonaws.com/dev/ingredients";
  return fetch(`${apiUrl}?q=${ingredientName}`)
    .then(response => response.json());
}

export default function IngredientSearch() {
  const [ingredients, setIngredients] = useState([]);
  const { dispatch } = useContext(StateContext);

  const onSearch = useCallback(ingredientName => {
    searchIngredients(ingredientName)
      .then(setIngredients)
      .catch(() => setIngredients([]));
  }, []);

  const onChange = useCallback(
    (ingredientIds) => {
      dispatch(actions.suggestRecipes(ingredientIds));
    },
    [dispatch]
  );

  return (
    <>
      <Select
        allowClear={true}
        mode="multiple"
        style={{ width: '100%' }}
        size="large"
        filterOption={false}
        placeholder="Search ingredients"
        onChange={onChange}
        notFoundContent={'No ingredients found'}
        onSearch={onSearch}
      >
        {ingredients.map(ingredient => <Option key={ingredient.id}>{ingredient.name}</Option>)}
      </Select>
    </>
  );
}
