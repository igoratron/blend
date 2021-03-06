package store

import "database/sql"

type MySQLCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type Ingredient struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (i Ingredient) Values() []interface{} {
	return []interface{}{i.Id, i.Name}
}

type Record interface {
	Values() []interface{}
}

type IngredientRecipe struct {
	IngredientId string
	RecipeId     string
}

func (i IngredientRecipe) Values() []interface{} {
	return []interface{}{i.IngredientId, i.RecipeId}
}

type Queryable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type Ingredients struct {
	Matching int      `json:"matching"`
	Total    int      `json:"total"`
	List     []string `json:"list"`
}

type Recipe struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Url         string      `json:"url"`
	Ingredients Ingredients `json:"ingredients"`
	ImagePath   string      `json:"imagePath"`
}

type ShoppingList struct {
	Ingredients []ShoppingListItem `json:"ingredients"`
}

type ShoppingListItem struct {
	Name   string  `json:"name"`
	Amount float32 `json:"amount"`
	Unit   string  `json:"unit"`
}
