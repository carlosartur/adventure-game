package models

import (
	"database/sql"
	"fmt"
	"log"

	"adventure-game/database"
)

type Player struct {
	Model
	Name            string
	Hability        int
	Luck            int
	Energy          int
	Provisions      int
	InitialHability int
	InitialLuck     int
	InitialEnergy   int
	Paragraph       int
	Items           []Item
}

func (Player) GetTableName() string {
	return `player`
}

func (p Player) Create() Player {
	p.Id = p.Model.Create(
		p.GetTableName(),
		p.Name,
		p.Paragraph,
		p.Hability,
		p.Luck,
		p.Energy,
		p.InitialHability,
		p.InitialLuck,
		p.InitialEnergy,
		p.Provisions,
	)

	return p
}

func (p Player) Update() Player {

	params := map[string]interface{}{
		"name":             p.Name,
		"paragraph_id":   	p.Paragraph,
		"hability":   	    p.Hability,
		"luck":   	        p.Luck,
		"energy":   	    p.Energy,
		"initial_hability": p.InitialHability,
		"initial_luck":   	p.InitialLuck,
		"initial_energy":   p.InitialEnergy,
		"provisions":   	p.Provisions,
	}

	p.Model.Update(
		p.GetTableName(),
		p.Id,
		params,
	)

	return p
}

func (p Player) Retrieve() []Player {
	response, err := database.ExecSql(`SELECT * FROM `+p.GetTableName()+` WHERE name LIKE ?;`, "%"+p.Name+"%")
	
	defer database.CloseDB()

	if err != nil {
		log.Fatal(err)
		return nil
	}

	list, err := p.BuildResponse(response)

	if err != nil {
		fmt.Println("Erro ao obter item.")
		return nil
	}

	return list
}

func (p Player) RetrieveOneById() Player {
	response, err := database.ExecSql(`SELECT 
		id, 
		name, 
		hability, 
		luck, 
		energy, 
		provisions, 
		initial_hability,
		initial_luck,
		initial_energy,
		paragraph_id
	FROM `+p.GetTableName()+` WHERE id = ?;`, p.Id)
	
	defer database.CloseDB()

	if err != nil {
		log.Fatal(err)
		return Player{}
	}

	list, err := p.BuildResponse(response)

	if err != nil || len(list) < 1 {
		fmt.Println("Erro ao obter item.")
		return Player{}
	}

	return list[0]
}

func (p Player) Delete() {
	p.Model.Delete(p.GetTableName(), p.Id)
}

func (p Player) BuildResponse(rows *sql.Rows) ([]Player, error) {
	var response []Player

	for rows.Next() {
		var newPlayer Player

		err := rows.Scan(
			&newPlayer.Id,
			&newPlayer.Name,
			&newPlayer.Hability,
			&newPlayer.Luck,
			&newPlayer.Energy,
			&newPlayer.Provisions,
			&newPlayer.InitialHability,
			&newPlayer.InitialLuck,
			&newPlayer.InitialEnergy,
			&newPlayer.Paragraph,
		)

		if err != nil {
			log.Fatal(err)

			return nil, err
		}

		response = append(response, newPlayer)
	}

	return response, nil
}
