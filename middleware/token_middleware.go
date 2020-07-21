package middleware

import (
	"github.com/gin-gonic/gin"
	"go-rest/app"
	"go-rest/app/jwt"
	tokenService "go-rest/app/token"
	"go-rest/app/utils"
	"net/http"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := utils.ExtractToken(context)
		tokenValidationError := tokenService.VerifyToken(token)
		if tokenValidationError != nil {
			context.JSON(http.StatusUnauthorized, app.ErrJson{Err: tokenValidationError.Error()})
			context.Abort()
			return
		}

		accessTokenClaims, claimsErr := tokenService.GetToken(context)
		if claimsErr != nil {
			context.JSON(http.StatusUnauthorized, app.ErrJson{Err: claimsErr.Error()})
			context.Abort()
			return
		}

		refreshToken, refreshTokenSearchErr := tokenService.GetTokenByUUID(accessTokenClaims.Uuid)
		if refreshTokenSearchErr != nil {
			context.JSON(http.StatusUnauthorized, app.ErrJson{Err: jwt.TokenExpired})
			context.Abort()
			return
		}

		refreshTokenValidationError := tokenService.VerifyToken(refreshToken.RefreshToken)
		if refreshTokenValidationError != nil {
			context.JSON(http.StatusUnauthorized, app.ErrJson{Err: jwt.TokenExpired})
			context.Abort()
			return
		}

		if refreshToken.UserName != accessTokenClaims.Subject {
			context.JSON(http.StatusUnauthorized, app.ErrJson{Err: app.ErrAccessDenied})
			context.Abort()
			return
		}

		context.Next()
	}
}
