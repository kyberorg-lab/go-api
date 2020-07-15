package user

import (
	"fmt"
	"go-rest/app"
	"go-rest/app/crypto"
	"go-rest/app/database"
	"go-rest/app/database/model"
	"os"
)

type OldUser struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetSampleUser() OldUser {
	return OldUser{
		ID:       1,
		Username: "username",
		Password: "password",
	}
}

func CreateSuperUser() {
	var result string

	superUserName, nameFromEnvExists := os.LookupEnv(app.EnvSuperUserName)
	if !nameFromEnvExists {
		superUserName = app.DefaultSuperUserName
	}

	secretKeyPassword, secretKeyPasswordFromEnvExists := os.LookupEnv(app.EnvEncryptSecretKeyPassword)
	if !secretKeyPasswordFromEnvExists {
		secretKeyPassword = app.DefaultSecretKeyPassword
	}

	superUserPassword, passwordFromEnvExists := os.LookupEnv(app.EnvSuperUserPassword)
	var encryptedPassword string
	if passwordFromEnvExists {
		encPassword, err := crypto.EncryptString(superUserPassword, secretKeyPassword)
		if err != nil {
			fmt.Println("") //TODO handle
		}
		encryptedPassword = encPassword
	} else {
		encPassword, err := crypto.EncryptString(app.DefaultSuperUserPassword, secretKeyPassword)
		if err != nil {
			fmt.Println("") //TODO handle
		}
		encryptedPassword = encPassword
	}

	superUser, err := FindUserByName(superUserName)
	if err != nil {
		//we have to create it
		database.DBConn.Create(&model.User{
			Username: superUserName,
			Password: encryptedPassword,
			Scopes: []model.Scope{
				{Name: app.DefaultSuperUserScope},
			},
		})
		result = "created. Name: "
		result += superUserName
		result += " Password: value from "
		result += app.EnvSuperUserPassword
		result += " env var or default password '"
		result += app.DefaultSuperUserPassword
		result += "'"
	} else {
		//checking if password is up-to-date
		if superUser.Password != encryptedPassword {
			superUser.Password = encryptedPassword
			database.DBConn.Save(&superUser)
			result = "already exists (password updated)"
		} else {
			result = "skipping (already exists)"
		}
	}

	fmt.Println("Creating SuperUser", ".....", result)
}

func FindUserByName(username string) (model.User, error) {
	var user model.User
	result := database.DBConn.First(&user, "username = ?", username)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}
