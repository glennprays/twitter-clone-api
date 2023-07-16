package controller

import (
	"log"
	"net/http"
	"time"
	"twitter-clone-api/config/database"
	"twitter-clone-api/middleware"
	"twitter-clone-api/models"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func LoginUser(c *gin.Context) {
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	var login models.LoginRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := session.Run(c,
		"MATCH (u:User {username: $username, password: $password}) RETURN count(u) = 1",
		map[string]any{
			"username": login.Username,
			"password": login.Password,
		})
	if err != nil {
		log.Fatal(err)
		return
	}

	if result.Next(c) {
		count := result.Record().Values[0].(bool)
		log.Println(count)
		if count {
			middleware.CreateToken(c, login.Username, "user", 3600)

			c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
}

func LogoutAccount(c *gin.Context) {
	middleware.ResetUserToken(c)
	c.JSON(http.StatusOK, models.Response{
		Status:  200,
		Message: "Logout successful",
	})
}

func WhoAmI(c *gin.Context) {

	username, role, err := middleware.GetUsernameAndRoleFromCookie(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	responseData := models.BasicUser{
		Username: username,
		Role:     role,
	}
	c.JSON(http.StatusOK, responseData)
}

func GetUserHandler(c *gin.Context) {
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	users, err := getUser(session, c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, users)

}

func getUser(session neo4j.SessionWithContext, c *gin.Context) ([]models.User, error) {

	result, err := session.Run(c, "MATCH (p:User) RETURN p.name AS name, p.email AS email, p.username AS username, p.created AS created, p.password AS password", nil)
	if err != nil {
		return nil, err
	}

	var users []models.User

	for result.Next(c) {
		record := result.Record()

		nameVal, nameExists := record.Get("name")
		emailVal, emailExists := record.Get("email")
		usernameVal, usernameExists := record.Get("username")
		createdVal, createdExists := record.Get("created")

		if !(nameExists || emailExists || usernameExists || createdExists) {
			continue
		}

		name, _ := nameVal.(string)
		email, _ := emailVal.(string)
		username, _ := usernameVal.(string)
		created, _ := createdVal.(time.Time)

		user := models.User{
			Name:     &name,
			Email:    &email,
			Username: username,
			Created:  &created,
		}
		users = append(users, user)
	}

	return users, nil
}

func CreateUserHandler(c *gin.Context) {
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	var user models.User

	err = c.ShouldBindJSON(&user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	err = createUser(c, session, user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// // createUser creates a new user node in the Neo4j database.
func createUser(c *gin.Context, session neo4j.SessionWithContext, user models.User) error {

	_, err := session.Run(c, "CREATE (u:User {name: $name, username: $username, email: $email, password: $password, created: datetime()})", map[string]interface{}{
		"name":     user.Name,
		"username": user.Username,
		"email":    user.Email,
		"password": user.Password,
	})
	if err != nil {
		return err
	}

	return nil
}
