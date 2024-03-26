package database

import (
	"database/sql"

	// Driver for SQLite3 database
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	database, err := sql.Open("sqlite3", `../game.db`)
	treatError(err)
	db = database
}

func treatError(err error) {
	if err != nil {
		panic(err)
	}
}

func ExecSql(sqlQuery string, args ...any) (*sql.Rows, error) {
	db := GetDB()
	stmt, err := db.Prepare(sqlQuery)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(args...)

	return rows, err
}

func GetDB() *sql.DB {
	return db
}
