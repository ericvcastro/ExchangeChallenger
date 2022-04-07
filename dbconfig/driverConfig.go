package dbConfig

import "fmt"

type Article struct {
	ID    int
	Title string
	Body  []byte
}

const PostgresDriver = "postgres"

const User = "postgres"

const Host = "localhost"

const Port = "5432"

const Password = "postgres"

const DbName = "UserExchange"

const TableName = "article"

var DataSourceName = fmt.Sprintf("host=%s port=%s user=%s "+
	"password=%s dbname=%s sslmode=disable", Host, Port, User, Password, DbName)
