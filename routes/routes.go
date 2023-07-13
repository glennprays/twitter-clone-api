package routes

import (
	"fmt"
	"net/http"
	"twitter-clone-api/config/database"
	"twitter-clone-api/controller"
	"twitter-clone-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("data/files/:filename", controller.GetFile)
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			driver, err := database.ConnectDB()
			if err != nil {
				fmt.Printf("Failed to connect to Neo4j: %v", err)
			} else {
				fmt.Println("Success Connect to DB")
			}
			driver.Close(c)
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome to Twitter-clone",
			})
		})
		api.GET("/whoami", controller.WhoAmI)
		api.DELETE("/auth/logout", middleware.AuthMiddleware("user"), controller.LogoutAccount)
		user := api.Group("/user")
		{
			user.POST("/auth/login", controller.LoginUser)
			user.POST("/tweets", middleware.AuthMiddleware("user"), controller.PostTweet)
		}
	}

	return router
}
