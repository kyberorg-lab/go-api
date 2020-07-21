package model

import "github.com/jinzhu/gorm"

type Token struct {
	gorm.Model
	UserName     string
	UserAgent    string
	RefreshToken string
	RefreshUuid  string
	Expires      int64
	IssuedAt     int64
}
