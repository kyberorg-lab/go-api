package structure

import "github.com/dgrijalva/jwt-go"

type AppClaims struct {
	Authorized bool     `json:"auth,omitempty"`
	Uuid       string   `json:"uuid,omitempty"`
	Scopes     []string `json:"sco, omitempty"`
	jwt.StandardClaims
}
