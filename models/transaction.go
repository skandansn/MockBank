package models

import "github.com/jinzhu/gorm"

type Transaction struct {
	gorm.Model
	SenderAccount   string  `gorm:"size:255;not null;" json:"senderAccount" binding:"required"`
	ReceiverAccount string  `gorm:"size:255;not null;" json:"receiverAccount" binding:"required"`
	Amount          float64 `gorm:"size:255;not null;" json:"amount" binding:"required"`
	TransactionDate string  `gorm:"size:255;not null;" json:"transactionDate" binding:"required"`
	Message         string  `gorm:"size:255;not null;" json:"message"`
}

func CreateTransaction(t Transaction) (Transaction, error) {

	err := DB.Create(&t).Error
	if err != nil {
		return Transaction{}, err
	}

	return t, nil
}

func GetTransactionsForSenderAccount(accountNumber string, startDate string, endDate string) ([]Transaction, error) {

	var transactions []Transaction

	if err := DB.Where("sender_account = ?", accountNumber).Where("transaction_date BETWEEN ? AND ?", startDate, endDate).Order("transaction_date desc").Find(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}

func GetNLatestTransactionsForSenderAccount(accountNumber string, n uint, startDate string, endDate string) ([]Transaction, error) {

	var transactions []Transaction

	if err := DB.Where("sender_account = ? OR receiver_account = ?", accountNumber, accountNumber).Where("transaction_date BETWEEN ? AND ?", startDate, endDate).Order("transaction_date desc").Limit(n).Find(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}

func GetTransactionsForBankAccounts(accountNumbers []string, startDate string, endDate string) ([]Transaction, error) {

	var transactions []Transaction

	if err := DB.Where("sender_account IN (?) OR receiver_account IN (?)", accountNumbers, accountNumbers).Where("transaction_date BETWEEN ? AND ?", startDate, endDate).Order("transaction_date desc").Find(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}
