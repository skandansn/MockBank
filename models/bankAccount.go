package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type BankAccount struct {
	gorm.Model
	CustomerID     uint    `gorm:"not null" json:"customer_id" binding:"required"`
	AccountType    string  `gorm:"size:255;not null;" json:"accountType" binding:"required"`
	AccountNumber  string  `gorm:"size:255;not null;" json:"accountNumber" binding:"required"`
	AccountBalance float64 `gorm:"size:255;not null;" json:"accountBalance" binding:"required"`
}

func CreateBankAccount(account BankAccount) (BankAccount, error) {

	err := DB.Create(&account).Error
	if err != nil {
		return BankAccount{}, err
	}

	return account, nil
}

func GetBankAccountsByCustomerId(cid uint) ([]BankAccount, error) {

	var account []BankAccount

	if err := DB.Where("customer_id = ?", cid).Find(&account).Error; err != nil {
		return account, err
	}

	return account, nil
}

func CheckIfAccountBelongsToCustomerAndBalance(accountNumber string, cid uint, amount float64) (BankAccount, error) {

	var account BankAccount

	if err := DB.Where("account_number = ? AND customer_id = ? ", accountNumber, cid).First(&account).Error; err != nil {
		return account, err
	}

	if account.AccountBalance < amount {
		return account, errors.New("Insufficient Balance")
	}

	return account, nil
}

func CheckIfAccountBelongsToCustomer(accountNumber string, cid uint) (BankAccount, error) {

	var account BankAccount

	if err := DB.Where("account_number = ? AND customer_id = ? ", accountNumber, cid).First(&account).Error; err != nil {
		return account, err
	}

	return account, nil
}

func CheckIfAccountExists(accountNumber string) (BankAccount, error) {

	var account BankAccount

	if err := DB.Where("account_number = ?", accountNumber).First(&account).Error; err != nil {
		return account, err
	}

	return account, nil
}

func UpdateSenderAndReceiverAccountBalance(senderAccountNumber string, receiverAccountNumber string, amount float64) (BankAccount, error) {

	var senderAccount BankAccount
	var receiverAccount BankAccount

	if err := DB.Where("account_number = ?", senderAccountNumber).First(&senderAccount).Error; err != nil {
		return senderAccount, err
	}

	if err := DB.Where("account_number = ?", receiverAccountNumber).First(&receiverAccount).Error; err != nil {
		return receiverAccount, err
	}

	senderAccount.AccountBalance = senderAccount.AccountBalance - amount
	receiverAccount.AccountBalance = receiverAccount.AccountBalance + amount

	if err := DB.Save(&senderAccount).Error; err != nil {
		return senderAccount, err
	}

	if err := DB.Save(&receiverAccount).Error; err != nil {
		return receiverAccount, err
	}

	return senderAccount, nil
}
