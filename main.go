package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type walletType struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}
type investor struct {
	ID     string       `json:"id"`
	User   string       `json:"user"`
	Wallet []walletType `json:"wallet"`
}

var investorWallet = []investor{
	{
		ID: "1", User: "Ben",
		Wallet: []walletType{
			{Currency: "btc", Amount: 0.2},
			{Currency: "doge", Amount: 0.6},
		},
	},
	{
		ID: "2", User: "Eric",
		Wallet: []walletType{
			{Currency: "btc", Amount: 3.0},
			{Currency: "doge", Amount: 5000},
		},
	},
}

func getAllInvestor(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, investorWallet)
}

func findUser(nameUser string) (*investor, error) {
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
// }

func main() {
	route := gin.Default()
	route.GET("/investor", getAllInvestor)
	route.PATCH("/deposit", depositToken)
	route.PATCH("/withdraw", withdrawToken)
	route.Run("localhost:8080")
}
