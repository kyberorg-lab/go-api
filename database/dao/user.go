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

type UserDao struct {
	UserName       string
	HashedPassword string
	Scopes         []model.Scope
	Active         bool
}

func NewUserDao() UserDao {
	return UserDao{}
}

func (userDao *UserDao) CreateUser() error {
	err := userDao.checkUserData()
	if err != nil {
		return err
	}
	user := userDao.userStoreToUserModel()

	result := database.DBConn.Create(&user)
	result = database.DBConn.Save(&user)
	return result.Error
}

func (userDao *UserDao) UpdateUser() error {
	err := userDao.checkUserData()
	if err != nil {
		return err
	}
	user := userDao.userStoreToUserModel()

	result := database.DBConn.Save(&user)
	return result.Error
}

func (userDao *UserDao) CountUsers() (int, error) {
	numberOfUsers := 0
	result := database.DBConn.Model(&model.User{}).Count(&numberOfUsers)
	return numberOfUsers, result.Error
}

func (userDao *UserDao) FindUserByName() (model.User, error) {
	var user model.User
	var scopes model.Scope
	database.DBConn.Model(&user).Related(&scopes, "Scopes")
	result := database.DBConn.Preload("Scopes", &scopes).First(&user, "username = ?", userDao.UserName)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (userDao *UserDao) FindUsersByScope(scope model.Scope, onlyActiveUsers bool) ([]model.User, error) {
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

func (userDao *UserDao) checkUserData() error {
	if userDao.UserName == "" || len(userDao.UserName) < 3 || len(userDao.UserName) > 255 {
		return errors.New(UsernameNotValid)
	}
	if userDao.HashedPassword == "" || len(userDao.HashedPassword) > 72 {
		return errors.New(PasswordNotValid)
	}
	return nil
}

func (userDao *UserDao) userStoreToUserModel() model.User {
	return model.User{
		Username: userDao.UserName,
		Password: userDao.HashedPassword,
		Scopes:   userDao.Scopes,
		Active:   userDao.Active,
	}
}
