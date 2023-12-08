package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/skandansn/webDevBankBackend/utils/token"
	"golang.org/x/crypto/bcrypt"
	"html"
	"strings"
)

type User struct {
	gorm.Model
	Username   string `gorm:"size:255;not null;unique" json:"username"`
	Password   string `gorm:"size:255;not null;" json:"password"`
	Type       string `gorm:"size:255;not null;" json:"type"`
	UserTypeID uint   `gorm:"size:255;" json:"user_type_id"`
}

func GetUserByID(uid uint) (User, error) {

	var u User

	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func (u *User) PrepareGive() {
	u.Password = ""
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string) (string, error) {

	var err error

	u := User{}

	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.ID, u.Type, u.UserTypeID)

	if err != nil {
		return "", err
	}

	return token, nil

}

func (u *User) SaveUser() (*User, error) {

	var err error
	err = DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) UpdateUser() (*User, error) {

	var err error

	err = DB.Model(&u).Updates(User{Username: u.Username, Password: u.Password, Type: u.Type, UserTypeID: u.UserTypeID}).Error

	if err != nil {
		return &User{}, err
	}

	return u, nil

}

func (u *User) BeforeSave() error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil

}
