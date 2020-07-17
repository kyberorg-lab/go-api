package user

import (
	"errors"
	"fmt"
	"github.com/kyberorg/go-utils/osutils"
	"go-rest/app"
	"go-rest/app/database"
	"go-rest/app/database/model"
	"golang.org/x/crypto/bcrypt"
)

const bCryptCost = 14

func CreateFirstUser() error {
	var result string

	numOfUsers, countError := CountUsers()
	if countError != nil {
		return countError
	}
	if numOfUsers > 0 {
		result = "skipping (already exists)"
		fmt.Println("Creating FirstUser", ".....", result)
		return nil
	}

	firstUserName, _ := osutils.GetEnv(app.EnvFirstUserName, app.DefaultFirstUserName)
	firstUserPassword, passwordFromEnv := osutils.GetEnv(app.EnvFirstUserPassword, app.DefaultFirstUserPassword)

	firstUser, firstUserExists := FindUserByName(firstUserName)
	if firstUserExists != nil {
		//we have to create it

		hashedPassword, hasError := hashPassword(firstUserPassword)
		if hasError != nil {
			fmt.Println("Failed to hash password", hasError)
			return hasError
		}

		database.DBConn.Create(&model.User{
			Username: firstUserName,
			Password: hashedPassword,
			Scopes: []model.Scope{
				{Name: app.DefaultFirstUserScope},
			},
			Active: true,
		})

		result = "created. Name: "
		result += firstUserName

		result += " Password "
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
		isPasswordSame, compareError := CheckPasswordForUser(firstUser, firstUserPassword)
		if compareError != nil {
			fmt.Println("Failed to compare hashes", compareError)
			return compareError
		}

		if !isPasswordSame {
			//password update needed
			hashedPass, hashError := hashPassword(firstUserPassword)
			if hashError != nil {
				fmt.Println("Failed to hash password", hashError)
				return hashError
			}
			firstUser.Password = hashedPass
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

func CheckPasswordForUser(user model.User, passwordCandidate string) (bool, error) {
	if len(passwordCandidate) == 0 {
		return false, nil
	}
	if len(user.Password) == 0 {
		return false, errors.New("user has empty password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordCandidate))
	return err == nil, err
}

func CountUsers() (int, error) {
	numberOfUsers := 0
	result := database.DBConn.Model(&model.User{}).Count(&numberOfUsers)
	return numberOfUsers, result.Error
}

func hashPassword(plainTextPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bCryptCost)
	return string(bytes), err
}
