package main

import (
	"log"
	"twitter-clone-api/routes"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	// run the router
	r := routes.SetupRouter()
	r.Run(":8080")
}
