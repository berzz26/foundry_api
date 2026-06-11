package database

import (
	"log"
	"os"

	"foundry/pkg/database"
)

func main() {
	err := database.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()

	// start server
}