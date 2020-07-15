package api

import (
	"github.com/gin-gonic/gin"
	"go-rest/app/utils"
	"net/http"
)

func CreateUserEndpoint(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, utils.ErrorJson("Not implemented yet"))
}
