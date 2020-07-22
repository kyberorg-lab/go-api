package dao

import (
	"errors"
	"github.com/kyberorg/go-api/database"
	"github.com/kyberorg/go-api/database/model"
)

const (
	UsernameNotValid = "username cannot be empty or not valid"
	PasswordNotValid = "password is not valid"
)

type UserStore struct {
	UserName       string
	HashedPassword string
	Scopes         []model.Scope
	Active         bool
}

func NewUserStore() UserStore {
	return UserStore{}
}

func (us *UserStore) CreateUser() error {
	err := us.checkUserData()
	if err != nil {
		return err
	}
	user := us.userStoreToUserModel()

	result := database.DBConn.Create(&user)
	result = database.DBConn.Save(&user)
	return result.Error
}

func (us *UserStore) UpdateUser() error {
	err := us.checkUserData()
	if err != nil {
		return err
	}
	user := us.userStoreToUserModel()

	result := database.DBConn.Save(&user)
	return result.Error
}

func (us *UserStore) CountUsers() (int, error) {
	numberOfUsers := 0
	result := database.DBConn.Model(&model.User{}).Count(&numberOfUsers)
	return numberOfUsers, result.Error
}

func (us *UserStore) FindUserByName() (model.User, error) {
	var user model.User
	var scopes model.Scope
	database.DBConn.Model(&user).Related(&scopes, "Scopes")
	result := database.DBConn.Preload("Scopes", &scopes).First(&user, "username = ?", us.UserName)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (us *UserStore) FindUsersByScope(scope model.Scope, onlyActiveUsers bool) ([]model.User, error) {
	var allUsers []model.User
	var foundUsers []model.User
	var allScopes model.Scope
	database.DBConn.Model(&allUsers).Related(&allScopes, "Scopes")
	result := database.DBConn.Preload("Scopes", &allScopes).Find(&allUsers)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, user := range allUsers {
		if onlyActiveUsers && user.Active == false {
			break
		}
		for _, userScope := range user.Scopes {
			if userScope == scope {
				foundUsers = append(foundUsers, user)
			}
		}
	}
	return foundUsers, nil
}

func (us *UserStore) checkUserData() error {
	if us.UserName == "" || len(us.UserName) < 3 || len(us.UserName) > 255 {
		return errors.New(UsernameNotValid)
	}
	if us.HashedPassword == "" || len(us.HashedPassword) > 72 {
		return errors.New(PasswordNotValid)
	}
	return nil
}

func (us *UserStore) userStoreToUserModel() model.User {
	return model.User{
		Username: us.UserName,
		Password: us.HashedPassword,
		Scopes:   us.Scopes,
		Active:   us.Active,
	}
}
