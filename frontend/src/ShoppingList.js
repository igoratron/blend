import React, { useContext, useState } from "react";
import { Button, Drawer, List } from "antd";
import { ShoppingCartOutlined } from "@ant-design/icons";

import { StateContext } from "./App";

function RecipeItem(item) {
  return <List.Item>
    { item.name }
    </List.Item>
}

export default function ShoppingList() {
  const { state } = useContext(StateContext);

  const [isLoading, setLoading] = useState(false);
  const [shoppingList, setShoppingList] = useState([]);

  const [areRecipesVisible, setRecipesVisible] = useState(false);
  const openRecipes = () => setRecipesVisible(true);
  const closeRecipes = () => setRecipesVisible(false);

  const [isShoppingListVisible, setShoppingListVisible] = useState(false);
  const openShoppingList = () => setShoppingListVisible(true);
  const closeShoppingList = () => setShoppingListVisible(false);

  const requestShoppingList = () => {
    const apiUrl = "https://mfuqctb1me.execute-api.eu-west-1.amazonaws.com/dev/shopping-list?recipeIds=";
    const recipeIds = state.recipePlan.map(r => r.id).join(',')

    if(!recipeIds) {
      return;
    }

    setLoading(true);
    openShoppingList();
    fetch(apiUrl + recipeIds)
      .then(r => r.json())
      .then(({ ingredients }) => setShoppingList(ingredients))
      .finally(() => setLoading(false))
  };

  return (
    <>
      <Button type="text" onClick={openRecipes}>
        <ShoppingCartOutlined />
        Your Meal Plan
      </Button>
      <Drawer
        title="Your Meal Plan"
        placement="right"
        width={338}
        onClose={closeRecipes}
        visible={areRecipesVisible}
      >
        <List
          dataSource={state.recipePlan}
          renderItem={RecipeItem}
        />

        <Button type="text" onClick={requestShoppingList}>Show shopping list</Button>

        <Drawer
          title="Shopping List"
          placement="right"
          onClose={() => { closeRecipes(); closeShoppingList()}}
          visible={isShoppingListVisible}
        >
          <List
            size="small"
            loading={isLoading}
            dataSource={shoppingList}
            renderItem={item => <List.Item>{item.amount} {item.unit} {item.name}</List.Item>}
          />
        </Drawer>
      </Drawer>
    </>
  );
}
