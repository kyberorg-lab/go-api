package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kyberorg/go-api/global"
	"github.com/kyberorg/go-api/global/json"
	"github.com/kyberorg/go-api/global/utils"
	"net/http"
)

type Profile struct {
	Username    string   `json:"username"`
	Permissions []string `json:"permissions"`
}

func GetProfileEndpoint(context *gin.Context) {
	token := utils.ExtractToken(context)
	tokenClaims, _ := tokenService.ExtractClaimsFromToken(token)

	profile := Profile{
		Username:    tokenClaims.Subject,
		Permissions: tokenClaims.Scopes,
	}
	context.JSON(http.StatusOK, profile)
}

func GetMySessionsEndpoint(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, json.ErrJson{Err: global.ErrNotImplemented})
}
