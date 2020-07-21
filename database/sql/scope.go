package sql

import (
	"github.com/kyberorg/go-api/database"
	"github.com/kyberorg/go-api/database/model"
)

type ScopeStore struct {
	ScopeName string
}

func (ss *ScopeStore) CreateNewScope() error {
	result := database.DBConn.Create(&model.Scope{
		Name: ss.ScopeName,
	})
	return result.Error
}

func (ss *ScopeStore) FindScopeByName() (model.Scope, error) {
	var scope model.Scope
	result := database.DBConn.First(&scope, "name = ?", ss.ScopeName)

	if result.Error != nil {
		return scope, result.Error
	}
	return scope, nil
}
