package database

import (
	"database/sql"
	"os"

	// Driver for SQLite3 database
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	dbPath := `../game.db`

	// Verifica se o arquivo do banco de dados existe
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Se não existir, cria o arquivo
		file, err := os.Create(dbPath)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	// Abre a conexão com o banco de dados
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	db = database
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
