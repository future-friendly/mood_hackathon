package models

import (
	"errors"
	"github.com/future-friednly/mood/backend/util"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `validate:"required"`
	Email string `validate:"required,email"`
	PasswordHash string `validate:"required"`
}

func (u User) isUnique() bool {
	var check User
	if result := gormDB.Where(&User{Email: u.Email}).First(&check); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return true
		}
	}
	return false
}

func NewUser(name, email, password string) error {
	user := &User{Name: name, Email: email, PasswordHash: util.HashPassword(password)}
	if err := validate.Struct(user); err != nil {
		return err
	}
	if !user.isUnique() {
		return util.AlreadyExists{Model: "user", Key: "email"}
	}
	if result := gormDB.Create(user); result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUser(email string) (*User, error){
	var user User
	if result := gormDB.Where(&User{Email: email}).First(&user); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (u *User) CheckPassword(password string) error {
	if util.HashPassword(password) != u.PasswordHash {
		return util.WrongCredentials{Login: u.Email}
	}
	return nil
}