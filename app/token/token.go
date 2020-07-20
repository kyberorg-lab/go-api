package token

import (
	"errors"
	"fmt"
	"go-rest/app/database"
	"go-rest/app/database/model"
	"go-rest/app/jwt"
	"go-rest/app/token/details"
	"time"
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

func GetTokenForUserAgent(userAgent string) (model.Token, error) {
	var token model.Token
	timeNow := time.Now().Unix()
	result := database.DBConn.Where("user_agent = ? AND expires > ?", userAgent, timeNow).First(&token)
	if result.Error != nil {
		fmt.Println("No token found for userAgent", userAgent)
		return model.Token{}, result.Error
	} else {
		return token, nil
	}

}
