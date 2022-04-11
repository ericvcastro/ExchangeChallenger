package basicfuncs

import "database/sql"

func createValuesTableWallet(db *sql.DB, table string, token_id int, user_id int, amount float64) {
	addValuesToTable := `INSERT INTO ` + table + ` VALUES ($1, $2, $3)`
	_, err := db.Exec(addValuesToTable, token_id, token_id, amount)
	if err != nil {
		panic(err)
	}
}
