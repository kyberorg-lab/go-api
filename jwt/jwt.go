package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

const (
	accessSecretEnv = "ACCESS_SECRET"
)

func CreateToken(userId uint64) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(15 * time.Minute).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	token, err := at.SignedString([]byte(os.Getenv(accessSecretEnv)))
	if err != nil {
		return "", err
	}
	return token, nil
}
