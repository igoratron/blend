package main

import (
  _"github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

  "encoding/json"
  "fmt"
  "io/ioutil"
  "os"
)

var tableName = "hellofresh-recipes"

type HelloFreshResponse struct {
  Count int
  Items []map[string]interface{}
}

func loadRecipies(path string) HelloFreshResponse {
  raw, err := ioutil.ReadFile(path)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  var recipies HelloFreshResponse
  json.Unmarshal(raw, &recipies)
  return recipies
}

func main() {
  sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
  }))

  svc := dynamodb.New(sess)

  recipies := loadRecipies("./recipes/hellofresh-favourites-1.json")
  fmt.Println("loaded", recipies.Count, "recipies")

  for _, recipe := range recipies.Items {
    av, err := dynamodbattribute.MarshalMap(recipe)
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    fmt.Println("Putting", recipe["name"])

    _, err = svc.PutItem(&dynamodb.PutItemInput{
      TableName: &tableName,
      Item:      av,
    })

    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
  }
}
