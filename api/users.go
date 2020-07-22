package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kyberorg/go-api/global"
	"github.com/kyberorg/go-api/global/json"
	"net/http"
)

func CreateUserEndpoint(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, json.ErrJson{Err: global.ErrNotImplemented})
}
