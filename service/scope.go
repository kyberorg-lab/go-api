package service

import (
	"fmt"
	"github.com/kyberorg/go-api/app"
	"github.com/kyberorg/go-api/database/dao"
	"github.com/kyberorg/go-api/database/model"
)

var (
	applicationScopes = []string{app.ScopeSuperAdmin, app.ScopeUser}

	scopeStore = dao.NewScopeStore()
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
	scopeStore.ScopeName = scope
	return scopeStore.FindScopeByName()
}

func createScope(scopeName string) {
	var result string
	scopeStore.ScopeName = scopeName
	_, err := scopeStore.FindScopeByName()
	if err != nil {
		createError := scopeStore.CreateNewScope()
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
