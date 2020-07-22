package dao

import (
	"github.com/kyberorg/go-api/database"
	"github.com/kyberorg/go-api/database/model"
	"github.com/kyberorg/go-api/global/structure"
	"time"
)

type TokenDao struct {
}

func NewTokenDao() TokenDao {
	return TokenDao{}
}

func (tokenDao TokenDao) GetTokenForUserAgent(userAgent string) (model.Token, error) {
	var token model.Token
	timeNow := time.Now().Unix()
	result := database.DBConn.Where("user_agent = ? AND expires > ?", userAgent, timeNow).First(&token)
	if result.Error != nil {
		return model.Token{}, result.Error
	} else {
		return token, nil
	}
}

func (tokenDao TokenDao) GetTokenByUUID(tokenUuid string) (model.Token, error) {
	var token model.Token
	result := database.DBConn.Where("refresh_uuid = ?", tokenUuid).First(&token)
	if result.Error != nil {
		return model.Token{}, result.Error
	} else {
		return token, nil
	}
}

func (tokenDao TokenDao) DeleteToken(tokenUuid string) error {
	tokenToBeDeleted, searchError := tokenDao.GetTokenByUUID(tokenUuid)
	if searchError != nil {
		return searchError
	}
	result := database.DBConn.Unscoped().Delete(&tokenToBeDeleted)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (tokenDao TokenDao) SaveToken(tokenDetails *structure.TokenDetails) error {
	refreshToken := model.Token{
		UserName:     tokenDetails.Subject,
		UserAgent:    tokenDetails.UserAgent,
		RefreshToken: tokenDetails.RefreshToken,
		RefreshUuid:  tokenDetails.RefreshUuid,
		Expires:      tokenDetails.RefreshTokenExpires,
		IssuedAt:     tokenDetails.CreatedAt,
	}

	saveResult := database.DBConn.Create(&refreshToken)
	return saveResult.Error
}
