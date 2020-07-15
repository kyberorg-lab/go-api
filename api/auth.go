package api

import (
	"github.com/gin-gonic/gin"
	"go-rest/app/jwt"
	"go-rest/app/user"
	"go-rest/app/utils"
	"net/http"
)

func AuthEndpoint(context *gin.Context) {
	var u user.User
	var sampleUser user.User = user.GetSampleUser()

	if err := context.ShouldBindJSON(&u); err != nil {
		context.JSON(http.StatusUnprocessableEntity, utils.ErrorJson("Invalid json provided"))
		return
	}
	//compare the user from request, with defined one
	if sampleUser.Username != u.Username || sampleUser.Password != u.Password {
		context.JSON(http.StatusUnauthorized, utils.ErrorJson("Please provide valid login details"))
		return
	}
	token, err := jwt.CreateToken(u.ID)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, utils.ErrorJson(err.Error()))
		return
	}
	context.JSON(http.StatusOK, token)
}
