package main

import (
	"github.com/gin-gonic/gin"
	"go-rest/api"
	"go-rest/app/database"
	"go-rest/app/scope"
	"log"
)

var (
	router = gin.Default()
)

func init() {
	database.InitDatabase()
	scope.CreateSuperUserScope()
}

func main() {
	handleStaticResources()

	router.POST("/auth", api.AuthEndpoint)
	userApi := router.Group("/users")
	{
		userApi.POST("", api.CreateUserEndpoint)
	}

	scopeApi := router.Group("/scope")
	{
		scopeApi.POST("", api.CreateScopeEndpoint)
	}

	defer database.CloseDatabase()
	log.Fatal(router.Run(":8080"))
}

func handleStaticResources() {
	router.Static("/static", "./assets")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
}
