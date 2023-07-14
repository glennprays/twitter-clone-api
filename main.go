package main

import (
	"fmt"
	"log"
<<<<<<< Updated upstream
	"os"
	"twitter-clone-api/routes"
=======
	"twitter-clone-api-Copy/routes"
>>>>>>> Stashed changes

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dirPath := os.Getenv("FILES_LOCATION")
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
	driver, err := neo4j.NewDriver("neo4j://localhost:7687", neo4j.BasicAuth("neo4j", "12345678", ""))
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close()

	r := gin.Default()
	r.GET("/people", routes.GetPeopleHandler(driver))
	r.POST("/people", routes.CreateUserHandler(driver))
	r.Run(":8080")
	log.Println("Server started on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}
