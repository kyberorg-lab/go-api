package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kyberorg/go-api/app"
	"net/http"
)

func CreateUserEndpoint(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, app.ErrJson{Err: app.ErrNotImplemented})
}
