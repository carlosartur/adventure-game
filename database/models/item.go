package models

import (
	"database/sql"
	"fmt"
	"log"
	"adventure-game/database"
)

type Item struct {
	Model
	Name string
}

func (Item) GetTableName() string {
	return `item`
}

func (i Item) Create() Item {
	i.Id = i.Model.Create(i.GetTableName(), i.Name)
	return i
}

func (i Item) Update() Item {
	params := map[string]interface{}{
		"name": i.Name,
	}

	i.Model.Update("Item", i.Id, params)
	return i
}

func (i Item) Retrieve() []Item {
	response, err := database.ExecSql(`SELECT * FROM `+i.GetTableName()+` WHERE name LIKE ?;`, "%"+i.Name+"%")

	if err != nil {
		log.Fatal(err)
		return nil
	}

	err, list := i.BuildResponse(response)

	if err != nil {
		fmt.Println("Erro ao obter item.")
		return nil
	}
	defer database.CloseDB()

	return list
}

func (i Item) Delete() {
	i.Model.Delete(i.GetTableName(), i.Id)
	defer database.CloseDB()
}

func (i Item) BuildResponse(rows *sql.Rows) (error, []Item) {
	var response []Item

	for rows.Next() {
		var newItem Item
		if err := rows.Scan(&newItem.Id, &newItem.Name); err != nil {
			log.Fatal(err)

			return err, nil
		}

		response = append(response, newItem)
	}

	return nil, response
}
