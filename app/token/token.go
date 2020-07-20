package token

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-rest/app/database"
	"go-rest/app/database/model"
	"go-rest/app/jwt"
	"go-rest/app/token/details"
	"go-rest/app/utils"
)

func SaveToken(tokenDetails *details.TokenDetails) error {
	if tokenDetails == nil {
		return errors.New("token details are empty")
	}

	_, queryErr := getTokenByUUID(tokenDetails.RefreshUuid)
	if queryErr == nil {
		//token already exist - nothing to do
		return nil
	}

	claims, err := jwt.ParseToken(tokenDetails.RefreshToken)
	if err != nil {
		return err
	}

	refreshToken := model.Token{
		UserName:     claims.Subject,
		UserAgent:    tokenDetails.UserAgent,
		RefreshToken: tokenDetails.RefreshToken,
		RefreshUuid:  tokenDetails.RefreshUuid,
		Expires:      tokenDetails.RtExpires,
		IssuedAt:     tokenDetails.CreatedAt,
	}

	saveResult := database.DBConn.Create(&refreshToken)
	return saveResult.Error
}

func VerifyToken(token string) error {
	_, err := jwt.ParseToken(token)
	return err
}

func GetToken(context *gin.Context) (jwt.AppClaims, error) {
	token := utils.ExtractToken(context)
	claims, err := jwt.ParseToken(token)
	return claims, err
}

func getTokenByUUID(tokenUuid string) (model.Token, error) {
	var token model.Token
	result := database.DBConn.Where("refresh_uuid = ?", tokenUuid).First(&token)
	if result.Error != nil {
		return model.Token{}, result.Error
	} else {
		return token, nil
	}
}
