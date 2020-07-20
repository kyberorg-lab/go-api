package middleware

import (
	"github.com/gin-gonic/gin"
	tokenService "go-rest/app/token"
	"go-rest/app/utils"
	"net/http"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := utils.ExtractToken(context)
		tokenValidationError := tokenService.VerifyToken(token)
		if tokenValidationError != nil {
			context.JSON(http.StatusUnauthorized, utils.ErrorJson(tokenValidationError.Error()))
			context.Abort()
			return
		}
		context.Next()
	}
}
