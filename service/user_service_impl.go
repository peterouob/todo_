package service

import (
	"errors"
	"github.com/peterouob/todo_/db"
	"github.com/peterouob/todo_/model"
	"gorm.io/gorm"
)

func loginUser(user model.User) (int64, error) {
	if err := db.DB.Where("username=?", user.Username).Where("password=?", user.Password).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return -1, errors.New("error in login user")
	}
	return user.ID, nil
}

func registerUser(user model.User) error {
	if err := db.DB.Where("username=?", user.Username).Where("password=?", user.Password).Find(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.DB.Create(&user).Error; err != nil {
			return errors.New("error in register user")
		}
	} else {
		return errors.New("have the same user")
	}
	return nil
}
