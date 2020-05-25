package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	_ "github.com/go-sql-driver/mysql"
)

func getDBCredentials() (MySQLCredentials, error) {
	secretName := "blend/db/dev"
	credentials := MySQLCredentials{}

	svc := secretsmanager.New(session.New())
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(input)
	if err == nil {
		err = json.Unmarshal([]byte(*result.SecretString), &credentials)
	}

	if err != nil {
		return credentials, err
	}

	return credentials, nil
}

func makePlaceholderString(columns int, statements int) string {
	result := make([]string, statements)

	for i := 0; i < statements; i += 1 {
		placeholder := strings.Repeat("?, ", columns)
		result[i] = fmt.Sprintf("(%s)", strings.TrimSuffix(placeholder, ", "))
	}

	return strings.Join(result, ", ")
}

func columns(names ...string) *[]string {
  return &names
}

func connect() *sql.DB {
	fmt.Println("Getting DB credentials")
	credentials, err := getDBCredentials()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/blend", credentials.Username, credentials.Password, credentials.Host, credentials.Port)
	fmt.Println("Connecting to the DB")
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println("Connected successfuly")
	return db
}

func insertInto(db Queryable, table string, columns *[]string, records *[]Record) error {
	placeholderString := makePlaceholderString(len(*columns), len(*records))
	columnNames := strings.Join(*columns, ", ")

	sqlStatement := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s", table, columnNames, placeholderString)

	var allRecordValues []interface{}
	for _, record := range *records {
		allRecordValues = append(allRecordValues, record.Values()...)
	}

	fmt.Printf("Query: %s with %#v\n", sqlStatement, allRecordValues)
	_, err := db.Exec(sqlStatement, allRecordValues...)

	if err != nil {
    if strings.Contains(err.Error(), "Error 1062") {
      return DuplicateEntryError{ message: err.Error() }
    }

		return err
	}

	return nil
}

func selectFrom(db Queryable, columns *[]string, table string, where ...interface{}) *sql.Row {
  var whereValues []interface{}

  sqlStatement := fmt.Sprintf("SELECT %s FROM %s", strings.Join(*columns, ", "), table)

  if len(where) > 1 {
    whereClause := where[0].(string)
    whereValues = where[1:]
    sqlStatement = strings.Join([]string{sqlStatement, "WHERE", whereClause}, " ")
  }

	fmt.Printf("Query: %s with %#v\n", sqlStatement, whereValues)
  return db.QueryRow(sqlStatement, whereValues...)
}
