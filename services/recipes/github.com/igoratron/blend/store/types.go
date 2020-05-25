package store

import "database/sql"

type MySQLCredentials struct {
	Username string `json:username`
	Password string `json:password`
	Host     string `json:host`
	Port     int    `json:port`
}

type Ingredient struct {
	Id   string
	Name string
}

func (i Ingredient) Values() []interface{} {
	return []interface{}{i.Id, i.Name}
}

type Record interface {
	Values() []interface{}
}

type IngredientRecipe struct {
	IngredientId string
	RecipeId    string
}

func (i IngredientRecipe) Values() []interface{} {
	return []interface{}{i.IngredientId, i.RecipeId}
}

type Queryable interface {
  Exec(query string, args ...interface{}) (sql.Result, error)
  Query(query string, args ...interface{}) (*sql.Rows, error)
  QueryRow(query string, args ...interface{}) *sql.Row
}

type DuplicateEntryError struct {
  message string
}

func (err DuplicateEntryError) Error() string {
  return err.message
}
