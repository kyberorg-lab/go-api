package scope

import (
	"fmt"
	"go-rest/app"
	"go-rest/app/database"
	"go-rest/app/database/model"
)

var db = database.DBConn

func CreateSuperUserScope() {
	_, err := FindScopeByName(app.DefaultSuperUserScope)
	if err != nil {
		fmt.Println("Creating", app.DefaultSuperUserScope, "scope")
		db.Create(&model.Scope{
			Name: app.DefaultSuperUserScope,
		})
	}
}

func FindScopeByName(name string) (model.Scope, error) {
	var scope model.Scope
	result := db.Find(scope, name)

	if result.Error != nil {
		return scope, result.Error
	}
	return scope, nil
}
