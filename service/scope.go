package service

import (
	"fmt"
	"github.com/kyberorg/go-api/app"
	"github.com/kyberorg/go-api/database/dao"
	"github.com/kyberorg/go-api/database/model"
)

var (
	applicationScopes = []string{app.ScopeSuperAdmin, app.ScopeUser}

	scopeDao = dao.NewScopeDao()
)

type ScopeService struct {
}

func NewScopeService() ScopeService {
	return ScopeService{}
}

func (ss *ScopeService) CreateScopes() {
	for _, scopeName := range applicationScopes {
		createScope(scopeName)
	}
}

func (ss *ScopeService) FindScopeByName(scope string) (model.Scope, error) {
	scopeDao.ScopeName = scope
	return scopeDao.FindScopeByName()
}

func createScope(scopeName string) {
	var result string
	scopeDao.ScopeName = scopeName
	_, err := scopeDao.FindScopeByName()
	if err != nil {
		createError := scopeDao.CreateNewScope()
		if createError != nil {
			fmt.Println("Failed to create", scopeName, "Error: ", createError)
			result = "failed (create failed)"
		} else {
			result = "done"
		}
	} else {
		result = "skipped (already exists)"
	}
	fmt.Println("Creating", scopeName, "scope", ".....", result)
}
