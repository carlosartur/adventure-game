package models

import (
	"log"
	"database/sql"
	"adventure-game/database"
)

type Model struct {
	Id int
}

func (m *Model) Create(tableName string, params ...interface{}) int {
	query := "INSERT INTO " + tableName + " VALUES (NULL, "
	for i := 0; i < len(params); i++ {
		if i > 0 {
			query += ", "
		}
		query += "?"
	}
	query += ")"

	_, err := database.ExecSql(query, params...)
	defer database.CloseDB()

	if err != nil {
		log.Fatal("Erro ao tentar executar query de inserção:", err)
	}

	
	id := m.LastInsertId(tableName)
	return int(id)
}

func (m *Model) LastInsertId(tableName string) int {
	query := "SELECT seq FROM sqlite_sequence WHERE name = '"+tableName+"';"
	row, err := database.ExecSql(query)
	defer database.CloseDB()

	if !row.Next() || err != nil {
		return 0
	}

	var seq int
	row.Scan(&seq)
	return seq
}

func (m *Model) Update(tableName string, id int, columnValues map[string]interface{}) {
    query := "UPDATE " + tableName + " SET "
    params := make([]interface{}, 0)
    i := 0
    for column, value := range columnValues {
        if i > 0 {
            query += ", "
        }
        query += column + " = ?"
        params = append(params, value)
        i++
    }
    query += " WHERE id = ?"
    params = append(params, id)

    _, err := database.ExecSql(query, params...)

    if err != nil {
		log.Fatal("Erro ao tentar executar query de alteração:", err)
    }
	defer database.CloseDB()
}

func (m *Model) Delete(tableName string, id int) {
	query := "DELETE FROM " + tableName + " WHERE id = ?"
	_, err := database.GetDB().Exec(query, id)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Model) GetOneById(tableName string, id int) (*sql.Rows, error) {
	query := "SELECT * FROM " + tableName + " WHERE id = ?"
	res, err := database.ExecSql(query, id)
	defer database.CloseDB()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return res, nil
}
