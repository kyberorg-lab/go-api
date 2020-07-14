package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	router = gin.Default()
)

const (
	accessSecret = "ACCESS_SECRET"
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var user = User{
	ID:       1,
	Username: "username",
	Password: "password",
}

func main() {
	router.POST("/login", Login)
	log.Fatal(router.Run(":8080"))
}

func Login(context *gin.Context) {
	var u User
	if err := context.ShouldBindJSON(&u); err != nil {
		context.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from request, with defined one
	if user.Username != u.Username || user.Password != u.Password {
		context.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	token, err := CreateToken(user.ID)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	context.JSON(http.StatusOK, token)
}

func CreateToken(userId uint64) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(15 * time.Minute).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	token, err := at.SignedString([]byte(os.Getenv(accessSecret)))
	if err != nil {
		return "", err
	}
	return token, nil
}
