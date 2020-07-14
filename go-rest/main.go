package main

import (
	"github.com/gin-gonic/gin"
	"go-rest/jwt"
	"go-rest/user"
	"log"
	"net/http"
)

var (
	router = gin.Default()
)

func main() {
	router.POST("/login", Login)
	log.Fatal(router.Run(":8080"))
}

func Login(context *gin.Context) {
	var u user.User
	var sampleUser user.User = user.GetSampleUser()

	if err := context.ShouldBindJSON(&u); err != nil {
		context.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from request, with defined one
	if sampleUser.Username != u.Username || sampleUser.Password != u.Password {
		context.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	token, err := jwt.CreateToken(u.ID)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	context.JSON(http.StatusOK, token)
}
