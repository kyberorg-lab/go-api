package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/kyberorg/go-api/app"
	"github.com/kyberorg/go-api/app/token/details"
	tokenService "github.com/kyberorg/go-api/app/token/service"
	userService "github.com/kyberorg/go-api/app/user"
	"github.com/kyberorg/go-api/database/model"
	"github.com/satori/go.uuid"
	"os"
	"time"
)

const (
	EmptyTokenStringError = "got empty string instead of token"
	EmptyTokenError       = "got empty token"
	MalformedTokenError   = "malformed token"
	TokenExpired          = "token is either expired or not active yet"
	MalformedClaimsError  = "got malformed claims"
	GeneralError          = "something went wrong"
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

	//Creating Refresh Token
	//go and check if user agent already have session
	alreadyStoredToken, tokenSearchError := tokenService.GetTokenForUserAgent(userAgent)
	userAgentHasValidToken := tokenSearchError == nil

	if userAgentHasValidToken {
		//re-use token
		tokenDetails.RefreshTokenExpires = alreadyStoredToken.Expires
		tokenDetails.RefreshUuid = alreadyStoredToken.RefreshUuid
		tokenDetails.RefreshToken = alreadyStoredToken.RefreshToken
	} else {
		//new token
		tokenDetails.RefreshTokenExpires = time.Now().Add(app.LifetimeRefreshToken).Unix()
		tokenDetails.RefreshUuid = uuid.NewV4().String()

		refreshTokenClaims := AppClaims{
			Uuid: tokenDetails.RefreshUuid,
			StandardClaims: jwt.StandardClaims{
				Subject:   user.Username,
				ExpiresAt: tokenDetails.RefreshTokenExpires,
				IssuedAt:  tokenDetails.CreatedAt,
				NotBefore: tokenDetails.CreatedAt,
			},
		}
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshTokenClaims)
		tokenDetails.RefreshToken, tokenSearchError = refreshToken.SignedString(signingKey)
		if tokenSearchError != nil {
			return nil, tokenSearchError
		}
	}

	tokenDetails.AccessTokenExpires = time.Now().Add(app.LifetimeAccessToken).Unix()

	tokenDetails.UserAgent = userAgent

	var err error
	//Creating Access Token
	accessTokenClaims := AppClaims{
		Authorized: true,
		Uuid:       tokenDetails.RefreshUuid,
		Scopes:     userService.GetScopeNames(user),
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Username,
			ExpiresAt: tokenDetails.AccessTokenExpires,
			IssuedAt:  tokenDetails.CreatedAt,
			NotBefore: tokenDetails.CreatedAt,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
	tokenDetails.AccessToken, err = accessToken.SignedString(signingKey)
	if err != nil {
		return nil, err
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
			return AppClaims{}, errors.New(MalformedTokenError)
		}
	} else {
		return AppClaims{}, errors.New(GeneralError)
	}
}
