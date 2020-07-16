package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Scopes   []Scope `json:"scopes"`
	Active   bool    `json:"active"`
}
