package models

import (
	"database/sql"
	"fmt"
	"log"
	"x/database"
)

type Item struct {
	Model
	Name string
}

func (_ Item) GetTableName() string {
	return `item`
}

func (i Item) Create() Item {
	i.Id = i.Model.Create(i.GetTableName(), i.Name)
	return i
}

func (i Item) Update() Item {
	i.Model.Update("Item", i.Id, i.Name)
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

	return list
}

func (i Item) Delete() {
	i.Model.Delete(i.GetTableName(), i.Id)
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
