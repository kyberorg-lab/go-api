package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"go-rest/app"
	"go-rest/app/database/model"
	"go-rest/app/token/details"
	tokenService "go-rest/app/token/service"
	userService "go-rest/app/user"
	"os"
	"time"
)

const (
	EmptyTokenStringError = "got empty string instead of token"
	EmptyTokenError       = "parsed to empty token"
	MalformedTokenError   = "malformed token"
	TokenExpired          = "token is either expired or not active yet"
	MalformedClaimsError  = "got malformed claims"
	GeneralError          = "something wrong"
)

var signingKey = []byte(os.Getenv(app.EnvJwtSecret))

type AppClaims struct {
	Authorized bool     `json:"auth,omitempty"`
	Uuid       string   `json:"uuid,omitempty"`
	Scopes     []string `json:"sco, omitempty"`
	jwt.StandardClaims
}

//noinspection GoNilness
func CreateToken(user model.User, userAgent string) (*details.TokenDetails, error) {
	tokenDetails := &details.TokenDetails{}
	tokenDetails.CreatedAt = time.Now().Unix()

	tokenDetails.AtExpires = time.Now().Add(15 * time.Minute).Unix()
	tokenDetails.AccessUuid = uuid.NewV4().String()

	tokenDetails.UserAgent = userAgent

	var err error
	//Creating Access Token
	accessTokenClaims := AppClaims{
		Authorized: true,
		Uuid:       tokenDetails.AccessUuid,
		Scopes:     userService.GetScopeNames(user),
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Username,
			ExpiresAt: tokenDetails.AtExpires,
			IssuedAt:  tokenDetails.CreatedAt,
			NotBefore: tokenDetails.CreatedAt,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
	tokenDetails.AccessToken, err = accessToken.SignedString(signingKey)
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	//go and check if user agent already have session
	alreadyStoredToken, tokenSearchError := tokenService.GetTokenForUserAgent(userAgent)
	userAgentHasValidToken := tokenSearchError == nil

	if userAgentHasValidToken {
		//re-use token
		tokenDetails.RtExpires = alreadyStoredToken.Expires
		tokenDetails.RefreshUuid = alreadyStoredToken.RefreshUuid
		tokenDetails.RefreshToken = alreadyStoredToken.RefreshToken
	} else {
		//new token
		tokenDetails.RtExpires = time.Now().Add(24 * time.Hour).Unix()
		tokenDetails.RefreshUuid = uuid.NewV4().String()

		refreshTokenClaims := AppClaims{
			Uuid: tokenDetails.RefreshUuid,
			StandardClaims: jwt.StandardClaims{
				Subject:   user.Username,
				ExpiresAt: tokenDetails.RtExpires,
				IssuedAt:  tokenDetails.CreatedAt,
				NotBefore: tokenDetails.CreatedAt,
			},
		}
		rt := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshTokenClaims)
		tokenDetails.RefreshToken, tokenSearchError = rt.SignedString(signingKey)
		if tokenSearchError != nil {
			return nil, tokenSearchError
		}
	}

	return tokenDetails, nil
}

func ParseToken(tokenString string) (AppClaims, error) {
	if len(tokenString) == 0 {
		return AppClaims{}, errors.New(EmptyTokenStringError)
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &AppClaims{}, func(t *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if parsedToken == nil {
		return AppClaims{}, errors.New(EmptyTokenError)
	}

	if parsedToken.Valid {
		if claims, ok := parsedToken.Claims.(*AppClaims); ok {
			return *claims, nil
		} else {
			return AppClaims{}, errors.New(MalformedClaimsError)
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return AppClaims{}, errors.New(MalformedTokenError)
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return AppClaims{}, errors.New(TokenExpired)
		} else {
			return AppClaims{}, errors.New(GeneralError)
		}
	} else {
		return AppClaims{}, errors.New(GeneralError)
	}
}
