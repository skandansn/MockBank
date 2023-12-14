package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type Employee struct {
	gorm.Model
	FirstName string `gorm:"size:255;not null;" json:"firstName"`
	LastName  string `gorm:"size:255;not null;" json:"lastName"`

	UserID uint `gorm:"not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID"`

	UserName string `gorm:"size:255;not null;unique" json:"username"`
	Email    string `gorm:"size:255;not null;" json:"email" binding:"required,email"`
	Phone    string `gorm:"size:255;not null;" json:"phone" binding:"required"`
	Branch   string `gorm:"size:255;not null;" json:"branch" binding:"required"`

	Customers []Customer `gorm:"foreignKey:EmployeeID"`

	Role string `gorm:"size:255;not null;" json:"role"  binding:"required,min=2,max=50" validate:"isValidRole"`

	ManagerID uint      `gorm:"" json:"manager_id"`
	Manager   *Employee `gorm:"foreignKey:ManagerID"`
}

func GetAllEmployees() ([]Employee, error) {

	var e []Employee

	if err := DB.Find(&e).Error; err != nil {
		return e, errors.New("Employee not found!")
	}

	return e, nil
}

func GetEmployeeById(eid uint) (Employee, error) {

	var e Employee

	if err := DB.First(&e, eid).Error; err != nil {

		return e, errors.New("Employee not found!")
	}

	return e, nil

}

func GetEmployeeFirstNameAndLastNameById(eid uint) (string, string, error) {

	var e Employee

	if err := DB.First(&e, eid).Error; err != nil {

		return "", "", errors.New("Employee not found!")
	}

	return e.FirstName, e.LastName, nil
}

func SaveEmployee(employee Employee) (*Employee, error) {

	var err error
	err = DB.Create(&employee).Error
	if err != nil {
		return &Employee{}, err
	}
	return &employee, nil
}

func DeleteEmployee(employee uint) error {

	var err error
	err = DB.Delete(&Employee{}, employee).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateEmployee(employee Employee) (*Employee, error) {

	var err error
	err = DB.Save(&employee).Error
	if err != nil {
		return &Employee{}, err
	}
	return &employee, nil
}
