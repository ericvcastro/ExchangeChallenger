package main

import (
	"context"
	"database/sql"
	"errors"
	"exchange/ExchengeChalenger/dbconfig"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

var db *sql.DB
var err error

type WalletType struct {
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	PriceInDollar float64 `json:"priceInDollar"`
	PriceInEuro   float64 `json:"priceInEuro"`
	TimeRateUsed  string  `json:"timeRateUsed"`
	TotalEuros    float64 `json:"totalEuros"`
	TotalDollar   float64 `json:"totalDollar"`
}
type Investor struct {
	ID                     int          `json:"id"`
	User                   string       `json:"user"`
	Wallet                 []WalletType `json:"wallet"`
	TotalAllCurrencyEuro   float64      `json:"totalAllCurrencyEuro"`
	TotalAllCurrencyDollar float64      `json:totalAllCurrencyDollar`
}

var investorWallet = []Investor{
	{
		ID: 1, User: "Ben",
		Wallet: []WalletType{
			{
				Currency:      "btc",
				Amount:        0.2,
				PriceInDollar: 46285.90,
				PriceInEuro:   39438.60,
				TimeRateUsed:  "08/09/2021 23:34",
				TotalEuros:    9257.18,
				TotalDollar:   7887.72,
			},
			{
				Currency:      "doge",
				Amount:        0.6,
				PriceInDollar: 0.25447,
				PriceInEuro:   0.21942,
				TimeRateUsed:  "08/09/2021 23:34",
				TotalEuros:    0.152682,
				TotalDollar:   0.131652,
			},
		},
		TotalAllCurrencyEuro:   9257.332682,
		TotalAllCurrencyDollar: 7887.851652,
	},
	{
		ID: 2, User: "Eric",
		Wallet: []WalletType{
			{
				Currency:      "btc",
				Amount:        3.0,
				PriceInDollar: 46285.90,
				PriceInEuro:   39438.60,
				TimeRateUsed:  "08/09/2021 23:34",
				TotalEuros:    118315.80,
				TotalDollar:   138857.70,
			},
			{
				Currency:      "doge",
				Amount:        5000,
				PriceInDollar: 0.25447,
				PriceInEuro:   0.21942,
				TimeRateUsed:  "08/09/2021 23:34",
				TotalEuros:    1097.10,
				TotalDollar:   1272.35,
			},
		},
		TotalAllCurrencyEuro:   228025.80,
		TotalAllCurrencyDollar: 140130.05,
	},
}

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

// func balanceInvestor(c *gin.Context) {
// 	user, okUser := c.GetQuery("user")

// 	if !okUser {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing any query parameter."})
// 		return
// 	}

// 	investorAndYourWallet, err := findUser(user)

// 	if err != nil {
// 		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User Not found"})
// 		return
// 	}

// 	tmpEuro := 0
// 	tmpDollar := 0

// 	for index, elem := range investorAndYourWallet.Wallet {

// 	}

// }

func CreateDB() {
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
		query := "CREATE TABLE " + dbconfig.TableName + "(user_id int primary key, user_name text, id_wallet int)"
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()

		res, err := db.ExecContext(ctx, query)

		if err != nil {
			panic(err)
		}

		println("Table Created", res)
	}
}

func main() {

	CreateDB()

	// sqlSelect()

	route := gin.Default()
	route.GET("/investor", getAllInvestor)
	route.PATCH("/deposit", depositToken)
	route.PATCH("/withdraw", withdrawToken)
	route.Run("localhost:8080")
}
