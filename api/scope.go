package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateScopeEndpoint(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{})
}
