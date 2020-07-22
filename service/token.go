package service

import (
	"errors"
	"github.com/kyberorg/go-api/database/dao"
	"github.com/kyberorg/go-api/database/model"
	"github.com/kyberorg/go-api/global"
	"github.com/kyberorg/go-api/global/structure"
	uuid "github.com/satori/go.uuid"
	"time"
)

var (
	tokenDao = dao.NewTokenDao()

	jwtService = NewJwtService()
)

type TokenService struct {
}

func NewTokenService() TokenService {
	return TokenService{}
}

func (s TokenService) CreateTokens(user model.User, userAgent string) (structure.TokenDetails, error) {
	tokenDetails := structure.TokenDetails{}
	tokenDetails.CreatedAt = time.Now().Unix()
	tokenDetails.Subject = user.Username
	tokenDetails.UserAgent = userAgent

	//Creating Refresh Token
	tokenCreateProblem := s.createRefreshToken(&tokenDetails)
	if tokenCreateProblem != nil {
		return structure.TokenDetails{}, tokenCreateProblem
	}

	//Creating Access Token
	tokenCreateProblem = s.createAccessToken(&tokenDetails)
	if tokenCreateProblem != nil {
		return structure.TokenDetails{}, tokenCreateProblem
	}

	return tokenDetails, nil
}

func (s TokenService) SaveToken(tokenDetails *structure.TokenDetails) error {
	if tokenDetails == nil {
		return errors.New("token details are empty")
	}

	_, queryErr := tokenDao.GetTokenByUUID(tokenDetails.RefreshUuid)
	if queryErr == nil {
		//token already exist - nothing to do
		return nil
	}

	saveResult := tokenDao.SaveToken(tokenDetails)
	return saveResult
}

func (s TokenService) ExtractClaimsFromToken(token string) (structure.AppClaims, error) {
	claims, err := jwtService.ParseToken(token)
	return claims, err
}

func (s TokenService) VerifyToken(token string) error {
	_, err := jwtService.ParseToken(token)
	return err
}

func (s TokenService) GetTokenByUUID(uuidString string) (model.Token, error) {
	return tokenDao.GetTokenByUUID(uuidString)
}

func (s TokenService) DeleteToken(tokenUuid string) interface{} {
	return tokenDao.DeleteToken(tokenUuid)
}

func (s TokenService) createRefreshToken(tokenDetails *structure.TokenDetails) error {
	//go and check if user agent already have session
	alreadyStoredToken, tokenSearchError := tokenDao.GetTokenForUserAgent(tokenDetails.UserAgent)
	userAgentHasValidToken := tokenSearchError == nil
	if userAgentHasValidToken {
		//re-use token
		tokenDetails.RefreshTokenExpires = alreadyStoredToken.Expires
		tokenDetails.RefreshUuid = alreadyStoredToken.RefreshUuid
		tokenDetails.RefreshToken = alreadyStoredToken.RefreshToken
	} else {
		//new token
		tokenDetails.RefreshTokenExpires = time.Now().Add(global.LifetimeRefreshToken).Unix()
		tokenDetails.RefreshUuid = uuid.NewV4().String()
		tokenGenerationProblem := jwtService.GenerateRefreshToken(tokenDetails)
		if tokenGenerationProblem != nil {
			return tokenGenerationProblem
		}
	}
	return nil
}

func (s TokenService) createAccessToken(tokenDetails *structure.TokenDetails) error {
	tokenDetails.AccessTokenExpires = time.Now().Add(global.LifetimeAccessToken).Unix()
	tokenGenerationProblem := jwtService.GenerateAccessToken(tokenDetails)
	if tokenGenerationProblem != nil {
		return tokenGenerationProblem
	}
	return nil
}
