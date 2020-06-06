import React, { useContext, useEffect, useState, useCallback } from "react";
import { Card, Progress, List, Tooltip, Typography, message } from "antd";
import { LinkOutlined, PlusOutlined, MinusOutlined } from "@ant-design/icons";

import { StateContext } from "./App";
import * as actions from "./actions";

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

const AddToPlan = addRecipe => (
  <button className="btn-ghost" onClick={addRecipe}>
    <PlusOutlined />
  </button>
);

const RemoveFromPlan = removeRecipe => (
  <button className="btn-ghost" onClick={removeRecipe}>
    <MinusOutlined />
  </button>
);

function PlanControls({ recipe }) {
  const { state, dispatch } = useContext(StateContext);

  const addRecipe = useCallback(() => {
    dispatch(actions.recipeAddedToPlan(recipe));
    message.success("Recipe added to the shopping list");
  }, [dispatch, recipe]);

  const removeRecipe = useCallback(() => {
    dispatch(actions.recipeRemovedFromPlan(recipe));
    message.success("Recipe removed from the shopping list");
  }, [dispatch, recipe]);

  const isRecipeInPlan = state.recipePlan.find(r => r.name === recipe.name);

  return (
    <Tooltip title="Add to plan">
      {isRecipeInPlan ? RemoveFromPlan(removeRecipe) : AddToPlan(addRecipe)}
    </Tooltip>
  );
}

function Link({ href }) {
  return (
    <Tooltip title="Go to the recipe">
      <a href={href} target="_blank" rel="noopener noreferrer">
        <LinkOutlined />
      </a>
    </Tooltip>
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
      locale={{ emptyText: "No recipes" }}
      renderItem={recipe => {
        const ingredientMatch = calculateMatch(recipe.ingredients);
        const imageUrl =
          "https://img.hellofresh.com/f_auto,fl_lossy,q_auto,w_610/hellofresh_s3" +
          recipe.imagePath;

        return (
          <List.Item>
            <Card
              actions={[
                <PlanControls recipe={recipe} />,
                <Link href={recipe.url} />
              ]}
              cover={<img src={imageUrl} className="cover" alt="" />}
            >
              <Card.Meta title={recipe.name} />
              <Progress percent={ingredientMatch} size="small" />
              <Typography.Paragraph ellipsis={{ rows: 3, expandable: true }}>
                {recipe.ingredients.list.join(", ")}
              </Typography.Paragraph>
            </Card>
          </List.Item>
        );
      }}
    />
  );
}
