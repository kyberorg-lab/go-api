package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kyberorg/go-api/app"
	tokenService "github.com/kyberorg/go-api/app/token"
	"net/http"
)

type Profile struct {
	Username    string   `json:"username"`
	Permissions []string `json:"permissions"`
}

func GetProfileEndpoint(context *gin.Context) {
	tokenClaims, _ := tokenService.GetToken(context)

	profile := Profile{
		Username:    tokenClaims.Subject,
		Permissions: tokenClaims.Scopes,
	}
	context.JSON(http.StatusOK, profile)
}

func GetMySessionsEndpoint(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, app.ErrJson{Err: app.ErrNotImplemented})
}
