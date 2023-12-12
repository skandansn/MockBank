package customer

import (
	"github.com/skandansn/webDevBankBackend/entity/bankAccount"
	"github.com/skandansn/webDevBankBackend/entity/card"
)

type Customer struct {
	CustomerId  uint   `gorm:"primary_key;auto_increment" json:"customerId" binding:"required"`
	FirstName   string `gorm:"size:255;not null;" json:"firstName" binding:"required"`
	LastName    string `gorm:"size:255;not null;" json:"lastName" binding:"required"`
	UserName    string `gorm:"size:255;not null;unique" json:"username" binding:"required"`
	Email       string `gorm:"size:255;not null;" json:"email" binding:"required,email"`
	Phone       string `gorm:"size:255;not null;" json:"phone" binding:"required"`
	DateOfBirth string `gorm:"size:255;not null;" json:"dateOfBirth" binding:"required"`
	Address     string `gorm:"size:255;not null;" json:"address" binding:"required"`

	Accounts []bankAccount.Account `gorm:"foreignKey:CustomerId" json:"accounts"`
	Cards    []card.Card           `gorm:"foreignKey:CustomerId" json:"cards"`
}

type CreateCustomerAccountInput struct {
	AccountType    string  `gorm:"size:255;not null;" json:"accountType" binding:"required"`
	AccountBalance float64 `gorm:"size:255;not null;" json:"accountBalance" binding:"required"`
}

type CreateCustomerInput struct {
	CustomerId    uint                         `gorm:"primary_key;auto_increment" json:"customerId" binding:"required"`
	Accounts      []CreateCustomerAccountInput `gorm:"foreignKey:CustomerId" json:"accounts"`
	CardNetwork   string                       `gorm:"size:255;not null;" json:"cardNetwork" binding:"required"`
	AppointmentId uint                         `gorm:"size:255;not null;" json:"appointmentId" binding:"required"`
}
