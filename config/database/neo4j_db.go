package database

import (
	"context"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func ConnectDB() (neo4j.DriverWithContext, error) {
	ctx := context.Background()
	uri := "neo4j://" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	// driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	driver, err := neo4j.NewDriverWithContext(
		uri,
		neo4j.BasicAuth(username, password, ""))

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return driver, nil
}
