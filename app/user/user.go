package user

import (
	"fmt"
	"github.com/kyberorg/go-utils/crypto/aesgcm"
	"github.com/kyberorg/go-utils/osutils"
	"go-rest/app"
	"go-rest/app/database"
	"go-rest/app/database/model"
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

func CreateFirstUser() error {
	var result string

	firstUserName, _ := osutils.GetEnv(app.EnvFirstUserName, app.DefaultFirstUserName)
	firstUserPassword, passwordFromEnv := osutils.GetEnv(app.EnvFirstUserPassword, app.DefaultFirstUserPassword)
	secretKeyPassword, _ := osutils.GetEnv(app.EnvEncryptSecretKeyPassword, app.DefaultSecretKeyPassword)
	salt, _ := osutils.GetEnv(app.EnvEncryptSalt, app.DefaultSalt)

	firstUser, firstUserExists := FindUserByName(firstUserName)
	if firstUserExists != nil {
		//we have to create it
		encryptedPassword, encryptError := aesgcm.EncryptString(firstUserPassword, secretKeyPassword, salt)
		if encryptError != nil {
			fmt.Println("Failed to encrypt password", encryptError)
			return encryptError
		}

		database.DBConn.Create(&model.User{
			Username: firstUserName,
			Password: encryptedPassword,
			Scopes: []model.Scope{
				{Name: app.DefaultFirstUserScope},
			},
			Active: true,
		})

		result = "created. Name: "
		result += firstUserName

		result += "Password "
		if passwordFromEnv {
			result += "(value of "
			result += app.EnvFirstUserPassword
			result += ")"
		} else {
			result += "'"
			result += app.DefaultFirstUserPassword
			result += "'"
		}
	} else {
		//checking if password is up-to-date
		decryptedPassword, decryptionError := aesgcm.DecryptString(firstUser.Password, secretKeyPassword, salt)
		if decryptionError != nil {
			fmt.Println("Failed to decrypt password", decryptionError)
			return decryptionError
		}

		if firstUserPassword != decryptedPassword {
			//password update needed
			encryptedPass, encError := aesgcm.EncryptString(firstUserPassword, secretKeyPassword, salt)
			if encError != nil {
				fmt.Println("Failed to encrypt password", encError)
				return encError
			}
			firstUser.Password = encryptedPass
			database.DBConn.Save(&firstUser)
			result = "already exists (password updated)"
		} else {
			result = "skipping (already exists)"
		}
	}

	fmt.Println("Creating FirstUser", ".....", result)
	return nil
}

func FindUserByName(username string) (model.User, error) {
	var user model.User
	result := database.DBConn.First(&user, "username = ?", username)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func SuperAdminsInSystemExist() (bool, error) {
	var users []model.User
	result := database.DBConn.Find(&users)
	if result.Error != nil {
		return false, result.Error
	}

	for _, user := range users {
		for _, s := range user.Scopes {
			if s.Name == app.ScopeSuperAdmin && user.Active {
				return true, nil
			}
		}
	}
	return false, nil
}
