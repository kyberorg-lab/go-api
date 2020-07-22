package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kyberorg/go-api/api"
	"github.com/kyberorg/go-api/database"
	"github.com/kyberorg/go-api/middleware"
	"github.com/kyberorg/go-api/service"
	"log"
)

var (
	router = gin.Default()

	scopeService = service.NewScopeService()
	userService  = service.NewUserService()
)

func init() {
	database.InitDatabase()
	scopeService.CreateScopes()

	//TODO continue layout modifications from hereon...
	err := userService.CreateFirstUser()
	if err != nil {
		otherSuperAdminsExist, searchError := userService.SuperAdminsInSystemExist()
		if !otherSuperAdminsExist || searchError != nil {
			panic("Failed to create first user and there are no other admins exist")
		}
	}
}

func main() {
	handleStaticResources()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	authApi := router.Group("/auth")
	{
		authApi.POST("login", api.LoginEndpoint)
		authApi.POST("refresh_token", middleware.TokenAuthMiddleware(), api.RefreshTokenEndpoint)
		authApi.POST("logout", middleware.TokenAuthMiddleware(), api.LogoutEndpoint)
	}

	profileApi := router.Group("/profile", middleware.TokenAuthMiddleware())
	{
		profileApi.GET("", api.GetProfileEndpoint)
		profileApi.GET("/sessions", api.GetMySessionsEndpoint)
	}

	userApi := router.Group("/users", middleware.TokenAuthMiddleware())
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
