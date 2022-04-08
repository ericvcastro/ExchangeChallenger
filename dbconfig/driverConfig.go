package dbconfig

import "fmt"

type UserWallet struct {
	ID    int
	Title string
	Body  []byte
}

const PostgresDriver = "postgres"

const User = "userwallet"

const Host = "localhost"

const Port = "5432"

const Password = "postgres"

const DbName = "postgres"

const TableName = "UserWallet"

var DataSourceName = fmt.Sprintf("host=%s port=%s user=%s "+
	"password=%s dbname=%s sslmode=disable", Host, Port, User, Password, DbName)

// user: UserExchange
