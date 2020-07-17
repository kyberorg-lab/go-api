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
