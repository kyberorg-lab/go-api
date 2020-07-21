package scope

import (
	"fmt"
	"github.com/kyberorg/go-api/app"
	"github.com/kyberorg/go-api/app/database"
	"github.com/kyberorg/go-api/app/database/model"
)

var (
	applicationScopes = []string{app.ScopeSuperAdmin, app.ScopeUser}
)

func CreateScopes() {
	for _, scopeName := range applicationScopes {
		createScope(scopeName)
	}
}

func createScope(scopeName string) {
	var result string

	_, err := FindScopeByName(scopeName)
	if err != nil {
		database.DBConn.Create(&model.Scope{
			Name: scopeName,
		})
		result = "done"
	} else {
		result = "skipped (already exists)"
	}
	fmt.Println("Creating", scopeName, "scope", ".....", result)
}

func FindScopeByName(name string) (model.Scope, error) {
	var scope model.Scope
	result := database.DBConn.First(&scope, "name = ?", name)

	if result.Error != nil {
		return scope, result.Error
	}
	return scope, nil
}
