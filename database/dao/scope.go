package dao

import (
	"github.com/kyberorg/go-api/database"
	"github.com/kyberorg/go-api/database/model"
)

type ScopeDao struct {
	ScopeName string
}

func NewScopeDao() ScopeDao {
	return ScopeDao{}
}

func (scopeDao *ScopeDao) CreateNewScope() error {
	result := database.DBConn.Create(&model.Scope{
		Name: scopeDao.ScopeName,
	})
	return result.Error
}

func (scopeDao *ScopeDao) FindScopeByName() (model.Scope, error) {
	var scope model.Scope
	result := database.DBConn.First(&scope, "name = ?", scopeDao.ScopeName)

	if result.Error != nil {
		return scope, result.Error
	}
	return scope, nil
}
