package database

import (
	"fmt"
	"log"
)

type Migration struct {
	Number int
	Query  string
}

func RunMigrations() {
	db := GetDB()

	defer db.Close()

	for _, migration := range getMigrations() {
		var err error
		var count int

		fmt.Println("Running migration #", migration.Number)
		fmt.Println(23)

		if migration.Number == 0 {
			_, err = db.Exec(migration.Query)

			if err != nil {
				log.Printf("%q: %s\n", err, migration.Query)
				fmt.Println(30)

				return
			}
			fmt.Println(34)

			continue
		}

		err = db.QueryRow(
			fmt.Sprintf(`SELECT COUNT(*) FROM migrations WHERE number = '%d';`, migration.Number),
		).Scan(&count)

		if err != nil {
			fmt.Println(fmt.Sprintf(`%v`, err))
			return
		}

		res, err := db.Exec(`SELECT * FROM sqlite_master WHERE type='table'`)

		fmt.Println(res, err)

		if count > 0 {
			continue
		}

		_, err = db.Exec(migration.Query)

		fmt.Println(56)

		if err != nil {
			log.Printf("%q: %s\n", err, migration.Query)
			return
		}

		_, err = db.Exec(
			fmt.Sprintf(`INSERT INTO migrations (number) VALUES ('%d');`, migration.Number),
		)

		fmt.Println(66)

		if err != nil {
			log.Printf("%q: %s\n", err, migration.Query)
			return
		}
	}

}

func getMigrations() []Migration {
	return []Migration{
		{Number: 0, Query: `
			CREATE TABLE IF NOT EXISTS migrations (
				id     INTEGER PRIMARY KEY AUTOINCREMENT,
				number INTEGER
			);`},
		{Number: 1, Query: `
			CREATE TABLE IF NOT EXISTS item (
				id   INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT
			);`},
		{Number: 2, Query: `
			CREATE TABLE IF NOT EXISTS player (
				id         INTEGER PRIMARY KEY AUTOINCREMENT,
				name       TEXT,
				hability   INTEGER,
				luck       INTEGER,
				energy     INTEGER,
				provisions INTEGER
			);`},
		{Number: 3, Query: `
			CREATE TABLE IF NOT EXISTS player_items (
				player_id INTEGER,
				item_id   INTEGER,
				FOREIGN KEY(player_id) REFERENCES player(id),
				FOREIGN KEY(item_id) REFERENCES item(id),
				PRIMARY KEY(player_id, item_id)
			);`},
		{Number: 4, Query: `
			CREATE TABLE IF NOT EXISTS paragraph (
				id      INTEGER PRIMARY KEY AUTOINCREMENT,
				context TEXT
			);`},
		{Number: 5, Query: `
			CREATE TABLE IF NOT EXISTS paragraph_option (
				paragraph_origin              INTEGER,
				paragraph_detiny              INTEGER,
				FOREIGN KEY(paragraph_origin) REFERENCES paragraph(id),
				FOREIGN KEY(paragraph_detiny) REFERENCES paragraph(id),
				PRIMARY KEY(paragraph_origin, paragraph_detiny)
			);`},
	}
}
