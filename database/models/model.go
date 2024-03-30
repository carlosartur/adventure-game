package models

import (
	"log"

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
	if err != nil {
		log.Fatal("Erro ao tentar executar query de inserção:", err)
	}

	
	id := m.LastInsertId(tableName)
	return int(id)
}

func (m *Model) LastInsertId(tableName string) int {
	query := "SELECT seq FROM sqlite_sequence WHERE name = '"+tableName+"';"
	row, err := database.ExecSql(query)
	
	if !row.Next() || err != nil {
		return 0
	}

	var seq int
	row.Scan(&seq)
	return seq
}

func (m *Model) Update(tableName string, id int, params ...interface{}) {
	query := "UPDATE " + tableName + " SET "
	for i := 0; i < len(params); i++ {
		if i > 0 {
			query += ", "
		}
		query += "?"
	}
	query += " WHERE id = ?"
	params = append(params, id)

	_, err := database.ExecSql(query, params...)
	if err != nil {
		log.Fatal("Erro ao tentar executar query de alteração:", err)
	}
}

func (m *Model) Delete(tableName string, id int) {
	query := "DELETE FROM " + tableName + " WHERE id = ?"
	_, err := database.ExecSql(query, id)
	if err != nil {
		log.Fatal(err)
	}
}
