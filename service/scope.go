package service

import (
	"fmt"
	"github.com/kyberorg/go-api/app"
	"github.com/kyberorg/go-api/database/sql"
)

var (
	applicationScopes = []string{app.ScopeSuperAdmin, app.ScopeUser}
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

func createScope(scopeName string) {
	var result string

	_, err := sql.ScopeStore.FindScopeByName(scopeName)
	if err != nil {
		createError := sql.ScopeStore.CreateNewScope(scopeName)
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
