import React, { useReducer, useEffect, useState } from "react";
import { Layout, Row, Col } from 'antd';

import 'antd/dist/antd.css';
import "./App.css";

import logoSrc from './assets/logo.png';
import * as actions from "./actions";
import IngredientSearch from "./ingredientSearch";
import Loading from "./Loading";
import RecipeResults from './recipeResult';

const { Header, Content } = Layout;

export const StateContext = React.createContext();

const initialState = {
  ingredientIdQuery: []
};

const reducer = (state, action) => {
  switch (action.type) {
    case actions.RECIPES_REQUESTED:
      return {
        ...state,
        requestedIngredientIds: action.payload
      };
    default:
      console.log("Unknown action:", action);
      return state;
  }
};

function ping() {
  return fetch(
    "https://mfuqctb1me.execute-api.eu-west-1.amazonaws.com/dev/ingredients?q="
  )
    .then(r => r.ok)
    .catch(() => false);
}

async function checkDBIsUp() {
  const isDBUp = await ping();
  if (!isDBUp) {
    await checkDBIsUp();
  }

  return Promise.resolve();
}

function App() {
  const [state, dispatch] = useReducer(reducer, initialState);
  const [isDBUp, setDBUp] = useState(false);

  useEffect(() => {
    checkDBIsUp()
      .then(() => setDBUp(true))
      .catch(() => setDBUp(false));
  }, []);

  if (!isDBUp) {
    return <Loading />;
  }

  return (
    <StateContext.Provider value={{ state, dispatch }}>
      <Layout>
        <Header className="header">
          <img className="header_logo" src={logoSrc} alt="Blend" />
        </Header>
        <Content className="u-bg-white u-pt-16">
          <Row justify="center" gutter={[16, 16]}>
            <Col xs={20} sm={12}>
              <IngredientSearch />
            </Col>
          </Row>
          <Row justify="center" gutter={[16, 16]}>
            <Col xs={20}>
              <RecipeResults />
            </Col>
          </Row>
        </Content>
      </Layout>
    </StateContext.Provider>
  );
}

export default App;
