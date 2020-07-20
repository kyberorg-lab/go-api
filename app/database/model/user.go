package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string
	Password string
	Scopes   []Scope `gorm:"many2many:user_scopes"`
	Active   bool
}
