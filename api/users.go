package api

import (
	"github.com/gin-gonic/gin"
	"go-rest/app"
	"net/http"
)

func CreateUserEndpoint(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, app.ErrJson{Err: app.ErrNotImplemented})
}
