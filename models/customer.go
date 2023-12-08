package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type Customer struct {
	gorm.Model
	FirstName   string `gorm:"size:255;not null;" json:"firstName"`
	LastName    string `gorm:"size:255;not null;" json:"lastName"`
	UserID      uint   `gorm:"not null" json:"user_id"`
	User        User   `gorm:"foreignKey:UserID"`
	UserName    string `gorm:"size:255;not null;unique" json:"username"`
	Email       string `gorm:"size:255;not null;" json:"email" binding:"required, email"`
	Phone       string `gorm:"size:255;not null;" json:"phone" binding:"required"`
	DateOfBirth string `gorm:"size:255;not null;" json:"dateOfBirth" binding:"required"`
	Address     string `gorm:"size:255;not null;" json:"address" binding:"required"`
}

func GetCustomerByUserID(cid uint) (Customer, error) {

	var c Customer

	err := DB.Where("user_id = ?", cid).First(&c).Error
	if err != nil {
		return Customer{}, errors.New("customer not found")
	}

	return c, nil

}

func GetCustomerById(id uint) (Customer, error) {

	var c Customer

	err := DB.Where("id = ?", id).First(&c).Error
	if err != nil {
		return Customer{}, errors.New("customer not found")
	}

	return c, nil
}

func (c *Customer) SaveCustomer() (*Customer, error) {

	var err error
	err = DB.Create(&c).Error
	if err != nil {
		return &Customer{}, err
	}
	return c, nil
}

func UpdateCustomer(c Customer) (Customer, error) {

	err := DB.Save(&c).Error
	if err != nil {
		return Customer{}, err
	}

	return c, nil
}
