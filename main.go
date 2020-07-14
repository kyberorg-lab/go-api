package main

import (
	"github.com/gin-gonic/gin"
	"go-rest/api"
	"go-rest/app/database"
	"log"
)

var (
	router = gin.Default()
)

func init() {
	database.InitDatabase()
}

func main() {
	handleStaticResources()

	router.POST("/auth", api.AuthEndpoint)
	userApi := router.Group("/users")
	{
		userApi.POST("", api.CreateUserEndpoint)
	}

	defer database.CloseDatabase()
	log.Fatal(router.Run(":8080"))
}

func handleStaticResources() {
	router.Static("/static", "./assets")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
}
