import React, { useContext, useEffect, useState } from "react";
import { Card, Progress, List } from "antd";

import { StateContext } from "./App";

function calculateMatch(ingredients) {
  const value = ingredients.matching / ingredients.total;
  return Math.round(value * 100);
}

function fetchRecipes(ingredientIds) {
  if (!ingredientIds || !ingredientIds.length) {
    return Promise.reject();
  }

  const ingredientIdQuery = ingredientIds.join(",");
  const apiUrl =
    "https://mfuqctb1me.execute-api.eu-west-1.amazonaws.com/dev/recipes";

  return fetch(`${apiUrl}?ingredientIds=${ingredientIdQuery}`).then(response =>
    response.json()
  );
}

export default function RecipeResults() {
  const { state } = useContext(StateContext);
  const [recipes, setRecipes] = useState([]);
  const { requestedIngredientIds } = state;

  useEffect(() => {
    fetchRecipes(requestedIngredientIds)
      .then(setRecipes)
      .catch(() => setRecipes([]));
  }, [requestedIngredientIds, setRecipes]);

  return (
    <List
      grid={{ gutter: 16, xs: 1, sm: 2, md: 3, column: 4 }}
      dataSource={recipes}
      locale={({emptyText: "No recipes"})}
      renderItem={recipe => {
        const ingredientMatch = calculateMatch(recipe.ingredients);
        const imageUrl =
          "https://img.hellofresh.com/f_auto,fl_lossy,q_auto,w_610/hellofresh_s3" +
          recipe.imagePath;

        return (
          <List.Item>
            <a href={recipe.url} target="_blank" rel="noopener noreferrer">
              <Card
                hoverable={true}
                cover={<img src={imageUrl} alt="" />}
              >
                <Card.Meta title={recipe.name} />
                <Progress percent={ingredientMatch} size="small" />
              </Card>
            </a>
          </List.Item>
        );
      }}
    />
  );
}
