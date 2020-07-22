package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/kyberorg/go-api/global"
	"github.com/kyberorg/go-api/global/structure"
	"os"
)

var signingKey = []byte(os.Getenv(global.EnvJwtSecret))

var (
	userService = NewUserService()
)

type JwtService struct {
}

func NewJwtService() JwtService {
	return JwtService{}
}

func (jwtService *JwtService) GenerateRefreshToken(tokenDetails *structure.TokenDetails) error {
	refreshTokenClaims := structure.AppClaims{
		Uuid: tokenDetails.RefreshUuid,
		StandardClaims: jwt.StandardClaims{
			Subject:   tokenDetails.Subject,
			ExpiresAt: tokenDetails.RefreshTokenExpires,
			IssuedAt:  tokenDetails.CreatedAt,
			NotBefore: tokenDetails.CreatedAt,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshTokenClaims)

	var tokenSignError error
	tokenDetails.RefreshToken, tokenSignError = refreshToken.SignedString(signingKey)
	if tokenSignError != nil {
		return tokenSignError
	}
	return nil
}

func (jwtService *JwtService) GenerateAccessToken(tokenDetails *structure.TokenDetails) error {
	user, subjectSearchError := userService.FindUserByName(tokenDetails.Subject)
	if subjectSearchError != nil {
		return subjectSearchError
	}
	accessTokenClaims := structure.AppClaims{
		Authorized: true,
		Uuid:       tokenDetails.RefreshUuid,
		Scopes:     userService.GetUserScopesNames(user),
		StandardClaims: jwt.StandardClaims{
			Subject:   tokenDetails.Subject,
			ExpiresAt: tokenDetails.AccessTokenExpires,
			IssuedAt:  tokenDetails.CreatedAt,
			NotBefore: tokenDetails.CreatedAt,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
	var signError error
	tokenDetails.AccessToken, signError = accessToken.SignedString(signingKey)
	if signError != nil {
		return signError
	}

	return nil
}

func (jwtService JwtService) ParseToken(tokenString string) (structure.AppClaims, error) {
	if len(tokenString) == 0 {
		return structure.AppClaims{}, errors.New(global.ErrEmptyTokenString)
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &structure.AppClaims{}, func(t *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if parsedToken == nil {
		return structure.AppClaims{}, errors.New(global.ErrEmptyToken)
	}

	if parsedToken.Valid {
		if claims, ok := parsedToken.Claims.(*structure.AppClaims); ok {
			return *claims, nil
		} else {
			return structure.AppClaims{}, errors.New(global.ErrMalformedClaims)
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return structure.AppClaims{}, errors.New(global.ErrMalformedToken)
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return structure.AppClaims{}, errors.New(global.ErrTokenExpired)
		} else {
			return structure.AppClaims{}, errors.New(global.ErrMalformedToken)
		}
	} else {
		return structure.AppClaims{}, errors.New(global.ErrGeneralError)
	}
}
