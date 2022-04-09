package main

import (
	"database/sql"
	"errors"
	"exchange/ExchengeChalenger/dbconfig"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

var db *sql.DB
var err error

// var userDB = 0

type WalletType struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}
type Investor struct {
	ID     int          `json:"id"`
	User   string       `json:"user"`
	Wallet []WalletType `json:"wallet"`
}

var investorWallet = []Investor{
	{
		ID: 1, User: "Ben",
		Wallet: []WalletType{
			{
				Currency: "btc",
				Amount:   0.2,
			},
			{
				Currency: "doge",
				Amount:   0.6,
			},
		},
	},
	{
		ID: 2, User: "Eric",
		Wallet: []WalletType{
			{
				Currency: "btc",
				Amount:   3.0,
			},
			{
				Currency: "doge",
				Amount:   5000,
			},
		},
	},
}

var testInvest []Investor

func getAllInvestor(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, investorWallet)
}

func findUser(nameUser string) (*Investor, error) {
	for i, b := range investorWallet {
		if b.User == nameUser {
			return &investorWallet[i], nil
		}
	}

	return nil, errors.New("Investor Not Found")
}

func depositToken(c *gin.Context) {
	user, okUser := c.GetQuery("user")
	currency, okCurrency := c.GetQuery("currency")
	amount, okAmount := c.GetQuery("amount")

	if !okUser || !okCurrency || !okAmount {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing any query parameter."})
		return
	}

	investorAndYourWallet, err := findUser(user)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User Not found"})
		return
	}

	amountFloat, errChangeToFloat := strconv.ParseFloat(amount, 64)

	if errChangeToFloat != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Erro to change Amount to Float"})
		return
	}

	for index, elem := range investorAndYourWallet.Wallet {
		if elem.Currency == currency && amountFloat > 0 {
			investorAndYourWallet.Wallet[index].Amount += amountFloat
		}
	}
	c.IndentedJSON(http.StatusOK, investorAndYourWallet)
}

func withdrawToken(c *gin.Context) {
	user, okUser := c.GetQuery("user")
	currency, okCurrency := c.GetQuery("currency")
	amount, okAmount := c.GetQuery("amount")

	if !okUser || !okCurrency || !okAmount {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing any query parameter."})
		return
	}

	investorAndYourWallet, err := findUser(user)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User Not found"})
		return
	}

	amountFloat, errChangeToFloat := strconv.ParseFloat(amount, 64)

	if errChangeToFloat != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Erro to change Amount to Float"})
		return
	}

	for index, elem := range investorAndYourWallet.Wallet {
		if elem.Currency == currency && amountFloat > 0 {
			if investorAndYourWallet.Wallet[index].Amount > amountFloat {
				investorAndYourWallet.Wallet[index].Amount -= amountFloat
			} else {
				investorAndYourWallet.Wallet[index].Amount = 0
			}
		}
	}
	c.IndentedJSON(http.StatusOK, investorAndYourWallet)
}

func Deposit(c *gin.Context) {

}

func SelectDataDB(name string) {
	query := `select * from ` + dbconfig.TableName
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	for row.Next() {
		investor := dbconfig.UserWallet{}
		err := row.Scan(&investor.ID, &investor.User)
		if err != nil {
			panic(err)
		}
		fmt.Println(investor)
	}
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
	route.GET("/investor", getAllInvestor)
	route.PATCH("/deposit", depositToken)
	route.PATCH("/withdraw", withdrawToken)
	route.Run("localhost:8080")
}
