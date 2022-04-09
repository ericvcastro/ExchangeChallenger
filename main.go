package main

import (
	"database/sql"
	"exchange/ExchengeChalenger/dbconfig"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

var db *sql.DB
var err error

type Investor struct {
	ID        int    `json:"id"`
	User_name string `json:"user"`
}

type Wallet struct {
	Wallet_id int `json:"wallet_id"`
	User_id   int `json:"user_id"`
}

type Tokens struct {
	Token_id int    `json:"token_id"`
	Currency string `json:"currency"`
}

type TokenWallet struct {
	Token_id  int `json:"token_id"`
	Wallet_id int `json:"wallet_id"`
	amount    int `json:"amount"`
}

// func getAllInvestor(c *gin.Context) {
// 	c.IndentedJSON(http.StatusOK, investorWallet)
// }

// func findUser(nameUser string) (*Investor, error) {
// 	for i, b := range investorWallet {
// 		if b.User == nameUser {
// 			return &investorWallet[i], nil
// 		}
// 	}

// 	return nil, errors.New("Investor Not Found")
// }

// func depositToken(c *gin.Context) {
// 	user, okUser := c.GetQuery("user")
// 	currency, okCurrency := c.GetQuery("currency")
// 	amount, okAmount := c.GetQuery("amount")

// 	if !okUser || !okCurrency || !okAmount {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing any query parameter."})
// 		return
// 	}

// 	investorAndYourWallet, err := findUser(user)

// 	if err != nil {
// 		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User Not found"})
// 		return
// 	}

// 	amountFloat, errChangeToFloat := strconv.ParseFloat(amount, 64)

// 	if errChangeToFloat != nil {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Erro to change Amount to Float"})
// 		return
// 	}

// 	for index, elem := range investorAndYourWallet.Wallet {
// 		if elem.Currency == currency && amountFloat > 0 {
// 			investorAndYourWallet.Wallet[index].Amount += amountFloat
// 		}
// 	}
// 	c.IndentedJSON(http.StatusOK, investorAndYourWallet)
// }

// func withdrawToken(c *gin.Context) {
// 	user, okUser := c.GetQuery("user")
// 	currency, okCurrency := c.GetQuery("currency")
// 	amount, okAmount := c.GetQuery("amount")

// 	if !okUser || !okCurrency || !okAmount {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing any query parameter."})
// 		return
// 	}

// 	investorAndYourWallet, err := findUser(user)

// 	if err != nil {
// 		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User Not found"})
// 		return
// 	}

// 	amountFloat, errChangeToFloat := strconv.ParseFloat(amount, 64)

// 	if errChangeToFloat != nil {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Erro to change Amount to Float"})
// 		return
// 	}

// 	for index, elem := range investorAndYourWallet.Wallet {
// 		if elem.Currency == currency && amountFloat > 0 {
// 			if investorAndYourWallet.Wallet[index].Amount > amountFloat {
// 				investorAndYourWallet.Wallet[index].Amount -= amountFloat
// 			} else {
// 				investorAndYourWallet.Wallet[index].Amount = 0
// 			}
// 		}
// 	}
// 	c.IndentedJSON(http.StatusOK, investorAndYourWallet)
// }

func Deposit(c *gin.Context) {
	user, okUser := c.GetQuery("user")
	currency, okCurrency := c.GetQuery("currency")
	amount, okAmount := c.GetQuery("amount")

	if !okUser || !okCurrency || !okAmount {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing any query parameter."})
		return
	}
	// Recuperando o User no DB
	println(`currency: ` + currency)
	println(`amount: ` + amount)
	username := SelectUserToTable(user, dbconfig.TableName)
	println(username.ID)
}

func SelectUserToTable(user string, table string) Investor {
	query := `select * from ` + table
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer row.Close()
	// mainUser := []Investor{}

	investor := Investor{}
	for row.Next() {
		err := row.Scan(&investor.ID, &investor.User_name)
		if err != nil {
			panic(err)
		}
		if investor.User_name == user {
			break
		}
	}
	return investor
}

func SelectDataDB(name string) {
	query := `select * from ` + dbconfig.TableName
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	investor := dbconfig.UserWallet{}
	for row.Next() {
		err := row.Scan(&investor.ID, &investor.User)
		if err != nil {
			panic(err)
		}
		// if investor.User == name {
		// 	return (investor)
		// }
	}
	fmt.Println(investor)
}

func UpdateDB() {
	query := `update ` + dbconfig.TableName + ` SET user_name = 'Eric' WHERE user_id = 2`
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer row.Close()
}

func main() {
	println("Acessando ", dbconfig.DbName)

	db, err = sql.Open(dbconfig.PostgresDriver, dbconfig.DataSourceName)

	if err != nil {
		panic(err.Error())
	} else {
		println("Connected!")
		println("")
	}

	defer db.Close()

	_, table_check := db.Query("select * from " + dbconfig.TableName + ";")
	if table_check == nil {
		println("Table is there")
	} else {
		dbconfig.CreateDB(db)
	}

	SelectDataDB("ben")

	route := gin.Default()
	// route.GET("/investor", getAllInvestor)
	route.PATCH("/deposit", Deposit)
	// route.PATCH("/withdraw", withdrawToken)
	route.Run("localhost:8080")
}
