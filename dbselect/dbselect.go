package dbselect

import "database/sql"

type Investor struct {
	ID        int    `json:"id"`
	User_name string `json:"user"`
}

type Wallet struct {
	Wallet_id int `json:"wallet_id"`
	User_id   int `json:"user_id"`
}

type TokenWallet struct {
	Token_id  int     `json:"token_id"`
	Wallet_id int     `json:"wallet_id"`
	Amount    float64 `json:"amount"`
}

type Tokens struct {
	Token_id    int     `json:"token_id"`
	Currency    string  `json:"currency"`
	PriceDollar float64 `json:"price_dollar"`
	PriceEuro   float64 `json:"price_euro"`
	TimeRate    string  `json:"time_rate"`
}

type HistoryType struct {
	User_id      int    `json:"user_id"`
	Transaction  string `json:"transaction"`
	TimeRealized string `json:"time_realized"`
}

func SelectUserToTable(db *sql.DB, user string, table string) Investor {
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

func SelectWalletToTable(db *sql.DB, user_Id int, table string) Wallet {
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

func SelectAllTokens(db *sql.DB, wallet_id int, table string) []TokenWallet {
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

func SelectAmountOfTable(db *sql.DB, token_id int, wallet_id int, table string) TokenWallet {
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

func SelectCurrencyToTable(db *sql.DB, token_id int, table string) Tokens {
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

func SelectTokenIdToTable(db *sql.DB, token_name string, table string) Tokens {
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

func SelectAllHistoryUser(db *sql.DB, user_id int, table string) []HistoryType {
	query := `select * from ` + table
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	userHistory := make([]HistoryType, 0)
	history := HistoryType{}
	for row.Next() {
		err := row.Scan(&history.User_id, &history.Transaction, &history.TimeRealized)
		if err != nil {
			panic(err)
		}

		if history.User_id == user_id {
			userHistory = append(userHistory, history)
		}
	}
	return userHistory
}
