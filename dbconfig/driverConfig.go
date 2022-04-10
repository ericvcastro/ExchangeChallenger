package dbconfig

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UserWallet struct {
	ID   int
	User string
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

func CreateDB(db *sql.DB) {
	query := "CREATE TABLE " + TableName + `(user_id int primary key not null, user_name text);`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := db.ExecContext(ctx, query)

	if err != nil {
		panic(err)
	}
	println("Table Created", res)

	addValuesToTable := `INSERT INTO ` + TableName + ` VALUES ($1, $2)`
	_, err = db.Exec(addValuesToTable, 1, "ben")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(addValuesToTable, 2, "eric")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(addValuesToTable, 3, "evandro")
	if err != nil {
		panic(err)
	}
	query = `CREATE TABLE Wallet (wallet_id int primary key not null);`
	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err = db.ExecContext(ctx, query)

	if err != nil {
		panic(err)
	}
	println("Table Wallet Created", res)

	queryWallet := `ALTER TABLE Wallet ADD COLUMN user_id INTEGER REFERENCES ` + TableName + ` (user_id);`
	_, err = db.Exec(queryWallet)
	if err != nil {
		panic(err)
	}

	addValuesToTable = `INSERT INTO Wallet VALUES ($1, $2)`
	_, err = db.Exec(addValuesToTable, 1, 1)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(addValuesToTable, 2, 2)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(addValuesToTable, 3, 3)
	if err != nil {
		panic(err)
	}

	query = `CREATE TABLE tokens (token_id int primary key not null, currency text, price_dollar float, price_euro float, time_rate text);`
	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err = db.ExecContext(ctx, query)

	if err != nil {
		panic(err)
	}
	println("Table tokens Created", res)

	addValuesToTable = `INSERT INTO tokens VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(addValuesToTable, 1, "btc", 42789.71, 39342.13, "08/04/2022 13:57")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(addValuesToTable, 2, "doge", 0.1534, 0.141, "08/04/2022 13:57")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(addValuesToTable, 3, "eth", 3273.28, 3009.50, "08/04/2022 13:57")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(addValuesToTable, 4, "ada", 1.05, 0.9616, "08/04/2022 13:57")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(addValuesToTable, 5, "xrp", 0.7621, 0.7007, "08/04/2022 13:57")
	if err != nil {
		panic(err)
	}
	query = `CREATE TABLE tokenWallet (token_id int REFERENCES tokens(token_id), wallet_id int REFERENCES wallet(wallet_id), amount float);`
	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err = db.ExecContext(ctx, query)

	if err != nil {
		panic(err)
	}
	println("Table tokenWallet Created", res)

	addValuesToTable = `INSERT INTO tokenWallet VALUES ($1, $2, $3)`
	_, err = db.Exec(addValuesToTable, 1, 1, 0.6)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(addValuesToTable, 2, 1, 3000.6)
	if err != nil {
		panic(err)
	}

	query = `CREATE TABLE history (user_id int REFERENCES userwallet(user_id), transaction text, time_realized text);`
	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err = db.ExecContext(ctx, query)

	if err != nil {
		panic(err)
	}
	println("Table history Created", res)
}
