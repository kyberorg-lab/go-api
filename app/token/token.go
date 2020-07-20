package token

import (
	"errors"
	"go-rest/app/database"
	"go-rest/app/database/model"
	"go-rest/app/jwt"
	"go-rest/app/token/details"
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

func VerifyAndExtractToken(token string) (jwt.AppClaims, error) {
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
