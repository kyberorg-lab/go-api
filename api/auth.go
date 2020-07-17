package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest/app/jwt"
	"go-rest/app/token"
	"go-rest/app/user"
	"go-rest/app/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginJson struct {
	Username string `json:"username" binding: "required,min=3,max=100"`
	Password string `json:"password" binding: "required,min=3,max=256"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func LoginEndpoint(context *gin.Context) {
	var loginJson LoginJson

	if err := context.ShouldBindJSON(&loginJson); err != nil {
		context.JSON(http.StatusUnprocessableEntity, utils.ErrorJsonWithError("Invalid json provided", err))
		return
	}

	foundUser, searchError := user.FindUserByName(loginJson.Username)
	if searchError != nil {
		fmt.Println("No such user ", loginJson.Username)
		context.JSON(http.StatusUnauthorized, utils.ErrorJson("Please provide valid login details"))
		return
	}

	isPasswordValid, compareError := user.CheckPasswordForUser(foundUser, loginJson.Password)
	if compareError != nil {
		if compareError == bcrypt.ErrMismatchedHashAndPassword {
			isPasswordValid = false
		} else {
			fmt.Println("Password hash compare error ", compareError)
			context.JSON(http.StatusInternalServerError, utils.ErrorJson("Hups something went wrong at our side"))
			return
		}
	}

	if !isPasswordValid {
		context.JSON(http.StatusUnauthorized, utils.ErrorJson("Please provide valid login details"))
		return
	}

	//user agent = UA header or ip
	userAgent := context.GetHeader("User-Agent")

	tokenDetails, err := jwt.CreateToken(foundUser, userAgent)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, utils.ErrorJson(err.Error()))
		return
	}

	_ = token.SaveToken(tokenDetails)

	context.JSON(http.StatusOK, Tokens{
		AccessToken:  tokenDetails.AccessToken,
		RefreshToken: tokenDetails.RefreshUuid,
	})
}

func RefreshTokenEndpoint(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, utils.ErrorJson("Not implemented yet"))
}

func LogoutEndpoint(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, utils.ErrorJson("Not implemented yet"))
}
