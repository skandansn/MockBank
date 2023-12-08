package bankAccount

type Account struct {
	AccountID      uint    `gorm:"primary_key;auto_increment" json:"accountId" binding:"required"`
	CustomerId     uint    `gorm:"primary_key;auto_increment" json:"customerId" binding:"required"`
	AccountType    string  `gorm:"size:255;not null;" json:"accountType" binding:"required"`
	AccountNumber  string  `gorm:"size:255;not null;" json:"accountNumber" binding:"required"`
	AccountBalance float64 `gorm:"size:255;not null;" json:"accountBalance" binding:"required"`
}
