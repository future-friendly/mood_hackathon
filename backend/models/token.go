package models

import (
	"crypto/rand"
	"encoding/hex"
	"gorm.io/gorm"
	"time"
)

type AuthToken struct {
	gorm.Model
	UserID uint
	Token string
	Expiry int64
}

func CreateAuthToken(id uint) (string, error){
	var check User
	if result := gormDB.First(&check, id); result.Error != nil {
		return "", result.Error
	}
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := AuthToken{
		UserID: id,
		Token: hex.EncodeToString(bytes),
		Expiry: time.Now().Add(time.Hour * time.Duration(24) * time.Duration(30)).Unix(),
	}
	if result := gormDB.Create(&token); result.Error != nil {
		return "", result.Error
	}
	return token.Token, nil
}

func GetToken(token string) (*AuthToken, error) {
	var authToken AuthToken
	if result := gormDB.Where(AuthToken{Token: token}).First(&authToken); result.Error != nil {
		return nil, result.Error
	}
	return &authToken, nil
}

func (t *AuthToken) GetUserFromToken() (*User, error) {
	var user User
	if result := gormDB.First(&user, t.UserID); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (t *AuthToken) IsExpired() bool {
	if time.Now().Unix() > t.Expiry {
		t.Delete()
		return true
	}
	return false
}

func (t *AuthToken) Delete() error {
	if result := gormDB.Delete(&t); result.Error != nil {
		return result.Error
	}
	return nil
}