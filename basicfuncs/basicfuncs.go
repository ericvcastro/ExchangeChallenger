package basicfuncs

import (
	"database/sql"
)

func CreateValuesTableWallet(db *sql.DB, table string, token_id int, user_id int, amount float64) {
	addValuesToTable := `INSERT INTO ` + table + ` VALUES ($1, $2, $3)`
	_, err := db.Exec(addValuesToTable, token_id, token_id, amount)
	if err != nil {
		panic(err)
	}
}

func AddValuesTable(db *sql.DB, table string, user_id int, transaction string, time_realized string) {
	addValuesToTable := `INSERT INTO ` + table + ` VALUES ($1, $2, $3)`
	_, err := db.Exec(addValuesToTable, user_id, transaction, time_realized)
	if err != nil {
		panic(err)
	}
}

func UpdateDB(db *sql.DB, table string, setCamp string, dataToUpdate string, whereModify string) {
	query := `update ` + table + ` SET ` + setCamp + ` = ` + dataToUpdate + ` WHERE ` + whereModify
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer row.Close()
}
