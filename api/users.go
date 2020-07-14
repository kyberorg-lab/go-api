package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUserEndpoint(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{})
}
