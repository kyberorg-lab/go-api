package api

import (
	"github.com/gin-gonic/gin"
	"go-rest/app/utils"
	"net/http"
)

func CreateScopeEndpoint(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, utils.ErrorJson("Not implemented yet"))
}
