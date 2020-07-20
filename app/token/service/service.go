package service

import (
	"go-rest/app/database"
	"go-rest/app/database/model"
	"time"
)

func GetTokenForUserAgent(userAgent string) (model.Token, error) {
	var token model.Token
	timeNow := time.Now().Unix()
	result := database.DBConn.Where("user_agent = ? AND expires > ?", userAgent, timeNow).First(&token)
	if result.Error != nil {
		return model.Token{}, result.Error
	} else {
		return token, nil
	}
}
