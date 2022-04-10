package main

import (
	"database/sql"
	"exchange/ExchengeChalenger/dbconfig"
	"fmt"
	"net/http"
	"strconv"

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
	Token_id    int     `json:"token_id"`
	Currency    string  `json:"currency"`
	PriceDollar float64 `json:"price_dollar"`
	PriceEuro   float64 `json:"price_euro"`
	TimeRate    string  `json:"time_rate"`
}

type TokenWallet struct {
	Token_id  int     `json:"token_id"`
	Wallet_id int     `json:"wallet_id"`
	Amount    float64 `json:"amount"`
}

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
	userId := SelectUserToTable(user, dbconfig.TableName)

	// Recuperando o Wallet ID do User no DB
	walletId := SelectWalletToTable(userId.ID, "wallet")

	// Recuperando o token ID no DB
	tokenId := SelectTokenIdToTable(currency, "tokens")

	// Recuperando o Montante do token do usuário
	amountSaved := SelectAmountOfTable(tokenId.Token_id, walletId.Wallet_id, "tokenwallet")

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
		UpdateDB("tokenwallet", "amount", amountstring, queryWhere)
		c.IndentedJSON(http.StatusOK, gin.H{"status": 200, "message": "Update amount sucess"})
	} else if amountSaved.Token_id == 0 && amountSaved.Wallet_id == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": 400, "message": "Bad Request"})
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
	userId := SelectUserToTable(user, dbconfig.TableName)

	// Recuperando o Wallet ID do User no DB
	walletId := SelectWalletToTable(userId.ID, "wallet")

	// Recuperando o token ID no DB
	tokenId := SelectTokenIdToTable(currency, "tokens")

	// Recuperando o Montante do token do usuário
	amountSaved := SelectAmountOfTable(tokenId.Token_id, walletId.Wallet_id, "tokenwallet")

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
		UpdateDB("tokenwallet", "amount", amountstring, queryWhere)

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
	userId := SelectUserToTable(user, dbconfig.TableName)

	// Recuperando o Wallet ID do User no DB
	walletId := SelectWalletToTable(userId.ID, "wallet")

	tokensAlls := SelectAllTokens(walletId.Wallet_id, "tokenwallet")

	balance := BalanceType{}
	balanceTotalEuro := 0.0
	balanceTotalDollar := 0.0

	for i, s := range tokensAlls {
		infoTokens := SelectCurrencyToTable(s.Token_id, "tokens")

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

func SelectAllTokens(wallet_id int, table string) []TokenWallet {
	query := `select * from ` + table
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	tokensWallet := make([]TokenWallet, 0)
	token := TokenWallet{}
	for row.Next() {
		err := row.Scan(&token.Token_id, &token.Wallet_id, &token.Amount)
		if err != nil {
			panic(err)
		}

		if token.Wallet_id == wallet_id {
			tokensWallet = append(tokensWallet, token)
		}
	}
	return tokensWallet
}

func SelectUserToTable(user string, table string) Investor {
	query := `select * from ` + table
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer row.Close()

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

func SelectWalletToTable(user_Id int, table string) Wallet {
	query := `select * from ` + table
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	wallet := Wallet{}
	for row.Next() {
		err := row.Scan(&wallet.Wallet_id, &wallet.User_id)
		if err != nil {
			panic(err)
		}
		if wallet.Wallet_id == user_Id {
			break
		}
	}
	return wallet
}

func SelectCurrencyToTable(token_id int, table string) Tokens {
	query := `select * from ` + table
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	token := Tokens{}
	for row.Next() {
		err := row.Scan(&token.Token_id, &token.Currency, &token.PriceDollar, &token.PriceEuro, &token.TimeRate)
		if err != nil {
			panic(err)
		}
		if token.Token_id == token_id {
			break
		}
	}
	return token
}

func SelectTokenIdToTable(token_name string, table string) Tokens {
	query := `select * from ` + table
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	token := Tokens{}
	for row.Next() {
		err := row.Scan(&token.Token_id, &token.Currency, &token.PriceDollar, &token.PriceEuro, &token.TimeRate)
		if err != nil {
			panic(err)
		}
		if token.Currency == token_name {
			break
		}
	}
	return token
}

func SelectAmountOfTable(token_id int, wallet_id int, table string) TokenWallet {
	query := `select * from ` + table
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	userHaveToken := false
	tokenSave := TokenWallet{}
	for row.Next() {
		err := row.Scan(&tokenSave.Token_id, &tokenSave.Wallet_id, &tokenSave.Amount)
		if err != nil {
			panic(err)
		}
		if tokenSave.Token_id == token_id && tokenSave.Wallet_id == wallet_id {
			userHaveToken = true
			break
		}
	}
	if userHaveToken {
		return tokenSave
	}
	emptyWallet := TokenWallet{}
	return emptyWallet
}

func UpdateDB(table string, setCamp string, dataToUpdate string, whereModify string) {
	query := `update ` + table + ` SET ` + setCamp + ` = ` + dataToUpdate + ` WHERE ` + whereModify
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

	route := gin.Default()
	// route.GET("/investor", getAllInvestor)
	route.PATCH("/deposit", Deposit)
	route.PATCH("/withdraw", Withdraw)
	route.GET("/balance", Balance)
	route.Run("localhost:8080")
}
