package main

import (
	"fmt"

	"github.com/dimitarvalkanov7/chaoscamp-demo/database"
	_ "github.com/lib/pq"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}
