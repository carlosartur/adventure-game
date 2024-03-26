package main

import (
	"fmt"

	"github.com/carlosartur/adventure-game/database/database"
)

// install
func init() {
	database.RunMigrations()
}

func main() {
	fmt.Println(`Vamos ver`)
}
