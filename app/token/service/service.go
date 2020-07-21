package service

import (
	"github.com/kyberorg/go-api/database"
	"github.com/kyberorg/go-api/database/model"
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
