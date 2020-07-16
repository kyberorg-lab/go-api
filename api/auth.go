package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest/app/database/model"
	"go-rest/app/jwt"
	"go-rest/app/user"
	"go-rest/app/utils"
	"net/http"
)

type AuthJson struct {
	Username string `json:"username" binding: "required,min=3,max=100"`
	Password string `json:"password" binding: "required,min=3,max=256"`
}

func AuthEndpoint(context *gin.Context) {
	var u model.User

	if err := context.ShouldBindJSON(&u); err != nil {
		context.JSON(http.StatusUnprocessableEntity, utils.ErrorJsonWithError("Invalid json provided", err))
		return
	}

	foundUser, searchError := user.FindUserByName(u.Username)
	if searchError != nil {
		fmt.Println("No such user ", u.Username)
		context.JSON(http.StatusUnauthorized, utils.ErrorJson("Please provide valid login details"))
		return
	}

	decryptedPassword, decryptError := user.DecryptPasswordForUser(foundUser)
	if decryptError != nil {
		fmt.Println("Decrypt error ", decryptError)
		context.JSON(http.StatusInternalServerError, utils.ErrorJson("Hups something went wrong at our side"))
		return
	}

	if u.Password != decryptedPassword {
		fmt.Println("No such user ", u.Username)
		context.JSON(http.StatusUnauthorized, utils.ErrorJson("Please provide valid login details"))
		return
	}

	token, err := jwt.CreateToken(foundUser)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, utils.ErrorJson(err.Error()))
		return
	}
	context.JSON(http.StatusOK, token)
}
