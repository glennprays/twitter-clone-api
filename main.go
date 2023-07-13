package main

import (
	"fmt"
	"log"
	"os"
	"twitter-clone-api/routes"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dirPath := "/app/data/files"
	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		// Directory doesn't exist, create it
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			fmt.Println("Failed to create directory:", err)
			return
		}
	}
}

func main() {

	// run the router
	r := routes.SetupRouter()
	r.Run(":8080")
}
