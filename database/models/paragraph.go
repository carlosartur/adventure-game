package models

import (
	"log"
	"adventure-game/utils"
	"adventure-game/database"
)

type Paragraph struct {
	Model
	Context string
	Options []Paragraph
}

func (Paragraph) GetTableName() string {
	return `paragraph`
}

func (p Paragraph) GetOneById(idStr string) Paragraph {
	id, _ := utils.ParseInt(idStr)
	rows, err := p.Model.GetOneById(p.GetTableName(), id)

	if err != nil {
		log.Fatal("Erro ao tentar buscar paragrafo")
	}

	rows.Next()
	
	err = rows.Scan(
		&p.Id,
		&p.Context,
	)

	if err != nil {
		log.Fatal(err)
		return Paragraph{}
	}

	p.Options = p.GetOptions()

	return p
}

func (p Paragraph) ValidateSelectedDestiny(idDestiny string) bool {
	id, _ := utils.ParseInt(idDestiny)
	for _, option := range p.Options {
		if option.Id == id {
			return true
		}
	}
	return false
}

func (p Paragraph) GetOptions() []Paragraph {
	var options []Paragraph
	
	rows, err := database.ExecSql(`SELECT
		id, context
	FROM
		paragraph_option po
		JOIN paragraph p ON po.paragraph_destiny = p.id
	WHERE
		po.paragraph_origin = ?`,
		p.Id,
	)

	if err != nil {
		log.Fatal(err)
		return options
	}

	for rows.Next() {
		paragraph_destiny := Paragraph{}

		err := rows.Scan(
			&paragraph_destiny.Id,
			&paragraph_destiny.Context,
		)

		if err != nil {
			log.Fatal(err)
			return options
		}

		options = append(options, paragraph_destiny)

	}
	return options
}