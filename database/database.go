package database

import (
	"database/sql"
	"os"
	"strings"

	// Driver for SQLite3 database
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	connect()
}

func connect() {
	dbPath := `./game.db`

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

	if strings.Contains(sqlQuery, "INSERT") || strings.Contains(sqlQuery, "UPDATE") {
		_, err := stmt.Exec(args...)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	
	rows, err := stmt.Query(args...)
	return rows, err
}

func GetDB() *sql.DB {
	if db == nil {
		connect()
		return db	
	}

	if err := db.Ping(); err != nil {
		if err := db.Close(); err != nil {
			panic(err)
		}
		connect()
	}

	return db
}

func CloseDB() {
	if db == nil {
		return
	}

	db.Close()
}
