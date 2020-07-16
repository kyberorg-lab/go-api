package scope

import (
	"fmt"
	"go-rest/app"
	"go-rest/app/database"
	"go-rest/app/database/model"
)

func CreateSuperUserScope() {
	var result string

	_, err := FindScopeByName(app.DefaultFirstUserScope)
	if err != nil {
		database.DBConn.Create(&model.Scope{
			Name: app.DefaultFirstUserScope,
		})
		result = "done"
	} else {
		result = "skipped (already exists)"
	}
	fmt.Println("Creating", app.DefaultFirstUserScope, "scope", ".....", result)
}

func FindScopeByName(name string) (model.Scope, error) {
	var scope model.Scope
	result := database.DBConn.First(&scope, "name = ?", name)

	if result.Error != nil {
		return scope, result.Error
	}
	return scope, nil
}
