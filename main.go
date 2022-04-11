package main

import (
	"database/sql"
	"exchange/ExchengeChalenger/dbconfig"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"

	"exchange/ExchengeChalenger/basicfuncs"
	"exchange/ExchengeChalenger/dbselect"

	"github.com/gin-gonic/gin"
)

var db *sql.DB
var err error

type PerCurrency struct {
	Currency               string  `json:"currency"`
	Amount                 float64 `json:"amount"`
	PriceInDollars         float64 `json:"price_in_dollars"`
	PriceInEuro            float64 `json:"price_in_euro"`
	TimeRateUsed           string  `json:"time_rate_used"`
	TotalEurosToCurrency   float64 `json:"total_Euros_to_currency"`
	TotalDollarsToCurrency float64 `json:"total_Dollars_to_currency"`
}

type BalanceType struct {
	PerCurrency  []PerCurrency `json:"per_currency"`
	TotalEuros   float64       `json:"total_Euros"`
	TotalDollars float64       `json:"total_Dollars"`
}

func Deposit(c *gin.Context) {
	user, okUser := c.GetQuery("user")
	currency, okCurrency := c.GetQuery("currency")
	amount, okAmount := c.GetQuery("amount")

	if !okUser || !okCurrency || !okAmount {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing any query parameter."})
		return
	}
	// Recuperando o User no DB
	userId := dbselect.SelectUserToTable(db, user, dbconfig.TableName)

	// Recuperando o Wallet ID do User no DB
	walletId := dbselect.SelectWalletToTable(db, userId.ID, "wallet")

	// Recuperando o token ID no DB
	tokenId := dbselect.SelectTokenIdToTable(db, currency, "tokens")

	// Recuperando o Montante do token do usuário
	amountSaved := dbselect.SelectAmountOfTable(db, tokenId.Token_id, walletId.Wallet_id, "tokenwallet")

	if amountSaved.Token_id != 0 && amountSaved.Wallet_id != 0 {
		amountFloat, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			panic(err)
		}

		newAmount := amountSaved.Amount + amountFloat

		amountstring := strconv.FormatFloat(newAmount, 'E', -1, 64)
		tokenIdString := strconv.Itoa(tokenId.Token_id)
		walletIdString := strconv.Itoa(walletId.Wallet_id)

		queryWhere := `token_id = ` + tokenIdString + ` AND wallet_id = ` + walletIdString
		basicfuncs.UpdateDB(db, "tokenwallet", "amount", amountstring, queryWhere)

		currentTime := time.Now()

		basicfuncs.AddValuesTable(db, "history", userId.ID, "deposit", currentTime.Format("2006-01-02 15:04:05"))
		c.IndentedJSON(http.StatusOK, gin.H{"status": 200, "message": "Update amount sucess"})
	} else if amountSaved.Token_id == 0 && amountSaved.Wallet_id == 0 {
		amountFloat, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			panic(err)
		}
		basicfuncs.CreateValuesTableWallet(db, "tokenwallet", tokenId.Token_id, walletId.Wallet_id, amountFloat)
	}
}

func Withdraw(c *gin.Context) {
	user, okUser := c.GetQuery("user")
	currency, okCurrency := c.GetQuery("currency")
	amount, okAmount := c.GetQuery("amount")

	if !okUser || !okCurrency || !okAmount {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing any query parameter."})
		return
	}
	// Recuperando o User no DB
	userId := dbselect.SelectUserToTable(db, user, dbconfig.TableName)

	// Recuperando o Wallet ID do User no DB
	walletId := dbselect.SelectWalletToTable(db, userId.ID, "wallet")

	// Recuperando o token ID no DB
	tokenId := dbselect.SelectTokenIdToTable(db, currency, "tokens")

	// Recuperando o Montante do token do usuário
	amountSaved := dbselect.SelectAmountOfTable(db, tokenId.Token_id, walletId.Wallet_id, "tokenwallet")

	if amountSaved.Token_id != 0 && amountSaved.Wallet_id != 0 {
		amountFloat, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			panic(err)
		}
		newAmount := 0.0
		if amountSaved.Amount > amountFloat {
			newAmount = amountSaved.Amount - amountFloat
		}
		amountstring := strconv.FormatFloat(newAmount, 'E', -1, 64)
		tokenIdString := strconv.Itoa(tokenId.Token_id)
		walletIdString := strconv.Itoa(walletId.Wallet_id)
		queryWhere := `token_id = ` + tokenIdString + ` AND wallet_id = ` + walletIdString
		basicfuncs.UpdateDB(db, "tokenwallet", "amount", amountstring, queryWhere)

		currentTime := time.Now()
		basicfuncs.AddValuesTable(db, "history", userId.ID, "withdraw", currentTime.Format("2006-01-02 15:04:05"))

		c.IndentedJSON(http.StatusOK, gin.H{"status": 200, "message": "Update amount sucess"})
	} else if amountSaved.Token_id == 0 && amountSaved.Wallet_id == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": 400, "message": "Bad Request"})
	}
}

func Balance(c *gin.Context) {
	user, okUser := c.GetQuery("user")
	if !okUser {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing any query parameter."})
		return
	}

	// Recuperando o User no DB
	userId := dbselect.SelectUserToTable(db, user, dbconfig.TableName)

	// Recuperando o Wallet ID do User no DB
	walletId := dbselect.SelectWalletToTable(db, userId.ID, "wallet")

	tokensAlls := dbselect.SelectAllTokens(db, walletId.Wallet_id, "tokenwallet")

	balance := BalanceType{}
	balanceTotalEuro := 0.0
	balanceTotalDollar := 0.0

	for i, s := range tokensAlls {
		infoTokens := dbselect.SelectCurrencyToTable(db, s.Token_id, "tokens")

		euroTotal := s.Amount * infoTokens.PriceEuro
		balanceTotalEuro += euroTotal

		dollarTotal := s.Amount * infoTokens.PriceDollar
		balanceTotalDollar += dollarTotal

		perCurrency := PerCurrency{
			Currency:               infoTokens.Currency,
			Amount:                 s.Amount,
			PriceInDollars:         infoTokens.PriceDollar,
			PriceInEuro:            infoTokens.PriceEuro,
			TimeRateUsed:           infoTokens.TimeRate,
			TotalEurosToCurrency:   euroTotal,
			TotalDollarsToCurrency: dollarTotal,
		}
		balance.PerCurrency = append(balance.PerCurrency, perCurrency)
		fmt.Println(i, perCurrency)
	}
	balance.TotalDollars = balanceTotalDollar
	balance.TotalEuros = balanceTotalEuro

	c.IndentedJSON(http.StatusOK, balance)
}

func History(c *gin.Context) {
	user, okUser := c.GetQuery("user")
	if !okUser {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing any query parameter."})
		return
	}

	userInfos := dbselect.SelectUserToTable(db, user, dbconfig.TableName)

	historyUser := dbselect.SelectAllHistoryUser(db, userInfos.ID, "history")

	c.IndentedJSON(http.StatusOK, historyUser)
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

	route := gin.Default()
	route.PATCH("/deposit", Deposit)
	route.PATCH("/withdraw", Withdraw)
	route.GET("/balance", Balance)
	route.GET("/history", History)
	route.Run("localhost:8080")
}
