package service

import (
	"errors"
	"fmt"
	"github.com/kyberorg/go-api/database/dao"
	"github.com/kyberorg/go-api/database/model"
	"github.com/kyberorg/go-api/global"
	"github.com/kyberorg/go-utils/osutils"
	"golang.org/x/crypto/bcrypt"
)

var (
	userDao = dao.NewUserDao()

	scopeService = NewScopeService()
)

const bCryptCost = 14

type UserService struct {
}

func NewUserService() UserService {
	return UserService{}
}

func (us *UserService) CreateFirstUser() error {
	var result string

	numOfUsers, countError := userDao.CountUsers()
	if countError != nil {
		return countError
	}
	if numOfUsers > 0 {
		result = "skipping (already exists)"
		fmt.Println("Creating FirstUser", ".....", result)
		return nil
	}

	firstUserName, _ := osutils.GetEnv(global.EnvFirstUserName, global.DefaultFirstUserName)
	firstUserPassword, passwordFromEnv := osutils.GetEnv(global.EnvFirstUserPassword, global.DefaultFirstUserPassword)

	firstUser, firstUserExists := us.FindUserByName(firstUserName)
	if firstUserExists != nil {
		//we have to create it
		hashedPassword, hasError := us.hashPassword(firstUserPassword)
		if hasError != nil {
			fmt.Println("Failed to hash password", hasError)
			return hasError
		}

		superAdminScope, err := scopeService.FindScopeByName(global.DefaultFirstUserScope)
		if err != nil {
			fmt.Println("There is no", global.DefaultFirstUserScope, "stored in Database")
			return err
		}

		userDao.UserName = firstUserPassword
		userDao.HashedPassword = hashedPassword
		userDao.Scopes = []model.Scope{superAdminScope}
		userDao.Active = true

		createError := userDao.CreateUser()

		if createError != nil {
			fmt.Println("First user cannot be created", createError)
			result = "failed"
		} else {
			result = "created. Name: "
			result += firstUserName

			result += " Password "
			if passwordFromEnv {
				result += "(value of "
				result += global.EnvFirstUserPassword
				result += ")"
			} else {
				result += "'"
				result += global.DefaultFirstUserPassword
				result += "'"
			}
		}
	} else {
		//checking if password is up-to-date
		isPasswordSame, compareError := us.CheckPasswordForUser(firstUser, firstUserPassword)
		if compareError != nil {
			fmt.Println("Failed to compare hashes", compareError)
			return compareError
		}

		if !isPasswordSame {
			//password update needed
			hashedPass, hashError := us.hashPassword(firstUserPassword)
			if hashError != nil {
				fmt.Println("Failed to hash password", hashError)
				return hashError
			}
			userDao.HashedPassword = hashedPass
			updateErr := userDao.UpdateUser()
			if updateErr != nil {
				fmt.Println("Password update error")
				result = "already exists (failed to update password)"
			} else {
				result = "already exists (password updated)"
			}
		} else {
			result = "skipping (already exists)"
		}
	}

	fmt.Println("Creating FirstUser", ".....", result)
	return nil
}

func (us *UserService) FindUserByName(name string) (model.User, error) {
	userDao.UserName = name
	return userDao.FindUserByName()
}

func (us *UserService) CheckPasswordForUser(user model.User, passwordCandidate string) (bool, error) {
	if len(passwordCandidate) == 0 {
		return false, nil
	}
	if len(user.Password) == 0 {
		return false, errors.New("user has empty password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordCandidate))
	return err == nil, err
}

func (us *UserService) SuperAdminsInSystemExist() (bool, error) {
	superAdminScopeName := global.ScopeSuperAdmin
	superAdminScope, scopeNotFound := scopeService.FindScopeByName(superAdminScopeName)
	if scopeNotFound != nil {
		return false, scopeNotFound
	}

	superAdmins, searchError := userDao.FindUsersByScope(superAdminScope, true)
	if searchError != nil {
		return false, searchError
	}

	superAdminsExist := len(superAdmins) > 0
	return superAdminsExist, nil
}

func (us *UserService) GetUserScopesNames(user model.User) []string {
	if user.Scopes == nil {
		return []string{}
	}
	var scopeNames []string
	scopes := user.Scopes
	for _, sco := range scopes {
		scopeNames = append(scopeNames, sco.Name)
	}
	return scopeNames
}

func (us *UserService) hashPassword(plainTextPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bCryptCost)
	return string(bytes), err
}
