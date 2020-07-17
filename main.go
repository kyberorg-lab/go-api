package main

import (
	"github.com/gin-gonic/gin"
	"go-rest/api"
	"go-rest/app/database"
	"go-rest/app/scope"
	"go-rest/app/user"
	"log"
)

var (
	router = gin.Default()
)

func init() {
	database.InitDatabase()
	scope.CreateScopes()

	err := user.CreateFirstUser()
	if err != nil {
		otherSuperAdminsExist, searchError := user.SuperAdminsInSystemExist()
		if !otherSuperAdminsExist || searchError != nil {
			panic("Failed to create first user and there are no other admins exist")
		}
	}
}

func main() {
	handleStaticResources()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	authApi := router.Group("/auth")
	{
		authApi.POST("login", api.LoginEndpoint)
		authApi.POST("refresh-token", api.RefreshTokenEndpoint)
		authApi.POST("logout", api.LogoutEndpoint)
	}

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
