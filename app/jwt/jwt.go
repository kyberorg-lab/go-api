package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"go-rest/app"
	"os"
	"time"
)

func CreateToken(userId uint64) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(15 * time.Minute).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	token, err := at.SignedString([]byte(os.Getenv(app.EnvJwtSecret)))
	if err != nil {
		return "", err
	}
	return token, nil
}
