package model

import "github.com/jinzhu/gorm"

type Scope struct {
	gorm.Model
	Name string
}
