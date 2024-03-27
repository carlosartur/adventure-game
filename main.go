package main

import (
	"fmt"

	"adventure-game/database"
)

// install
func init() {
	database.RunMigrations()
}

func main() {
	fmt.Println(`Vamos ver`)
}
