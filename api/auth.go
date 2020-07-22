package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kyberorg/go-api/app"
	"github.com/kyberorg/go-api/app/jwt"
	"github.com/kyberorg/go-api/app/token"
	tokenService "github.com/kyberorg/go-api/app/token"
	"github.com/kyberorg/go-api/app/utils"
	"github.com/kyberorg/go-api/service"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const (
	UserPassWrong         = "Please provide valid login details"
	CannotCreateTokens    = "Failed to create tokens"
	SuccessfullyLoggedOut = "Successfully logged out"
)

var (
	userService = service.NewUserService()
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

	foundUser, searchError := userService.FindUserByName(loginJson.Username)
	if searchError != nil {
		fmt.Println("No such user ", loginJson.Username)
		context.JSON(http.StatusUnauthorized, app.ErrJson{Err: UserPassWrong})
		return
	}

	isPasswordValid, compareError := userService.CheckPasswordForUser(foundUser, loginJson.Password)
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
	foundUser, userSearchError := userService.FindUserByName(refreshToken.UserName)
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
