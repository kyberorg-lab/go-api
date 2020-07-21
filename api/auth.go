package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest/app"
	"go-rest/app/jwt"
	"go-rest/app/token"
	tokenService "go-rest/app/token"
	"go-rest/app/user"
	"go-rest/app/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const (
	UserPassWrong         = "Please provide valid login details"
	CannotCreateTokens    = "Failed to create tokens"
	SuccessfullyLoggedOut = "Successfully logged out"
)

type LoginJson struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=3,max=256"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func LoginEndpoint(context *gin.Context) {
	var loginJson LoginJson

	if err := context.ShouldBindJSON(&loginJson); err != nil {
		context.JSON(http.StatusUnprocessableEntity, app.ErrJson{Err: app.ErrInvalidJson})
		return
	}

	foundUser, searchError := user.FindUserByName(loginJson.Username)
	if searchError != nil {
		fmt.Println("No such user ", loginJson.Username)
		context.JSON(http.StatusUnauthorized, app.ErrJson{Err: UserPassWrong})
		return
	}

	isPasswordValid, compareError := user.CheckPasswordForUser(foundUser, loginJson.Password)
	if compareError != nil {
		if compareError == bcrypt.ErrMismatchedHashAndPassword {
			isPasswordValid = false
		} else {
			fmt.Println("Password hash compare error ", compareError)
			context.JSON(http.StatusInternalServerError, app.ErrJson{Err: app.ErrGeneralError})
			return
		}
	}

	if !isPasswordValid {
		context.JSON(http.StatusUnauthorized, app.ErrJson{Err: UserPassWrong})
		return
	}

	//user agent = UA header plus ip
	userAgent := utils.GetUniqueUserAgent(context)

	tokenDetails, err := jwt.CreateToken(foundUser, userAgent)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, app.ErrJson{Err: CannotCreateTokens})
		return
	}

	_ = token.SaveToken(tokenDetails)

	context.JSON(http.StatusOK, Tokens{
		AccessToken:  tokenDetails.AccessToken,
		RefreshToken: tokenDetails.RefreshToken,
	})
}

func RefreshTokenEndpoint(context *gin.Context) {
	tokenClaims, _ := tokenService.GetToken(context)

	refreshToken, searchError := tokenService.GetTokenByUUID(tokenClaims.Uuid)
	if searchError != nil {
		context.JSON(http.StatusInternalServerError, app.ErrJson{Err: app.ErrGeneralError})
	}
	foundUser, userSearchError := user.FindUserByName(refreshToken.UserName)
	if userSearchError != nil {
		context.JSON(http.StatusForbidden, app.ErrJson{Err: app.ErrAccessDenied})
		return
	}

	newTokenPairDetails, tokenCreateError := jwt.CreateToken(foundUser, refreshToken.UserAgent)
	if tokenCreateError != nil {
		context.JSON(http.StatusUnprocessableEntity, app.ErrJson{Err: app.ErrAccessDenied})
		return
	}

	newTokens := Tokens{
		AccessToken:  newTokenPairDetails.AccessToken,
		RefreshToken: newTokenPairDetails.RefreshToken,
	}

	context.JSON(http.StatusCreated, newTokens)
}

func LogoutEndpoint(context *gin.Context) {
	tokenClaims, _ := tokenService.GetToken(context)
	delError := tokenService.DeleteToken(tokenClaims.Uuid)
	if delError != nil {
		context.JSON(http.StatusInternalServerError, app.ErrJson{Err: app.ErrGeneralError})
		return
	}
	context.JSON(http.StatusOK, app.MessageJson{Message: SuccessfullyLoggedOut})
}
