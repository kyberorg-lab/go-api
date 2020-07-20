package api

import (
	"github.com/gin-gonic/gin"
	"go-rest/app/jwt"
	tokenService "go-rest/app/token"
	"go-rest/app/utils"
	"net/http"
)

type Profile struct {
	Username    string   `json:"username"`
	Permissions []string `json: "permission"`
}

func GetProfileEndpoint(context *gin.Context) {
	token := utils.ExtractToken(context)
	tokenClaims, tokenValidationError := tokenService.VerifyAndExtractToken(token)
	if tokenValidationError != nil {
		switch tokenValidationError.Error() {
		case jwt.EmptyTokenStringError:
		case jwt.EmptyTokenError:
			context.JSON(http.StatusUnauthorized, utils.ErrorJson("Token is absent. This endpoint requires token"))
			return
		case jwt.TokenExpired:
			context.JSON(http.StatusUnauthorized, utils.ErrorJson(jwt.TokenExpired))
			return
		default:
			context.JSON(http.StatusUnauthorized, utils.ErrorJson("Token is malformed or unparseble. Create refresh token"))
			return
		}
	}

	profile := Profile{
		Username:    tokenClaims.Subject,
		Permissions: tokenClaims.Scopes,
	}
	context.JSON(http.StatusOK, profile)
}

func GetMySessionsEndpoint(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, utils.ErrorJson("Not implemented yet"))
}
