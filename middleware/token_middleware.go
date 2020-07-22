package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kyberorg/go-api/global"
	"github.com/kyberorg/go-api/global/json"
	"github.com/kyberorg/go-api/service"

	"github.com/kyberorg/go-api/global/utils"
	"net/http"
)

var (
	tokenService = service.NewTokenService()
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := utils.ExtractToken(context)
		tokenValidationError := tokenService.VerifyToken(token)
		if tokenValidationError != nil {
			context.JSON(http.StatusUnauthorized, json.ErrJson{Err: tokenValidationError.Error()})
			context.Abort()
			return
		}

		accessTokenClaims, claimsErr := tokenService.ExtractClaimsFromToken(token)
		if claimsErr != nil {
			context.JSON(http.StatusUnauthorized, json.ErrJson{Err: claimsErr.Error()})
			context.Abort()
			return
		}

		refreshToken, refreshTokenSearchErr := tokenService.GetTokenByUUID(accessTokenClaims.Uuid)
		if refreshTokenSearchErr != nil {
			context.JSON(http.StatusUnauthorized, json.ErrJson{Err: global.ErrTokenExpired})
			context.Abort()
			return
		}

		refreshTokenValidationError := tokenService.VerifyToken(refreshToken.RefreshToken)
		if refreshTokenValidationError != nil {
			context.JSON(http.StatusUnauthorized, json.ErrJson{Err: global.ErrTokenExpired})
			context.Abort()
			return
		}

		if refreshToken.UserName != accessTokenClaims.Subject {
			context.JSON(http.StatusUnauthorized, json.ErrJson{Err: global.ErrAccessDenied})
			context.Abort()
			return
		}

		context.Next()
	}
}
