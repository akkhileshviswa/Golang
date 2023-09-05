package model

import (
	"errors"
	"html"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"size:255;not null;unique" json:"username"`
	Password  string `gorm:"size:255;not null;" json:"-"`
	Groceries []Grocery
}

func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	user.Password = string(hashedPassword)

	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GetUserById(uid uint) (User, error) {
	var user User

	db, err := Setup()
	if err != nil {
		log.Println(err)
		return User{}, err
	}

	if err := db.Preload("Groceries").Where("id=?", uid).Find(&user).Error; err != nil {
		return user, errors.New("user not found")
	}

	return user, nil
}
