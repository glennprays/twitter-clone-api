package testing

import (
	"context"
	"log"
	"testing"
	"twitter-clone-api/config/database"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func TestConnectDB(t *testing.T) {
	ctx := context.Background()
	driver, err := database.ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	_, err = session.Run(ctx, "RETURN 1", nil)
	if err != nil {
		t.Fatalf("Failed to execute test query: %v", err)
	}
	if err := session.Close(ctx); err != nil {
		t.Fatalf("Failed to close Neo4j session: %v", err)
	}
	t.Log("Successfully connected to the Neo4j database.")
}
