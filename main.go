package main

import (
	"electronicsStore/database"
	"electronicsStore/router"
)

func main() {
	database.Connect()
	database.RunMigrations()

	router.Setup().Run(":8080")
}
