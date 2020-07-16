package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"go-rest/app"
	"go-rest/app/database/model"
	"os"
	"time"
)

func CreateToken(user model.User) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = user.ID
	atClaims["exp"] = time.Now().Add(15 * time.Minute).Unix()
	atClaims["scopes"] = user.Scopes

	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	token, err := at.SignedString([]byte(os.Getenv(app.EnvJwtSecret)))
	if err != nil {
		return "", err
	}
	return token, nil
}
