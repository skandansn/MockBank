package transactions

type Transaction struct {
	TransactionId   uint    `gorm:"primary_key;auto_increment" json:"transactionId" binding:"required"`
	SenderAccount   string  `gorm:"size:255;not null;" json:"senderAccount" binding:"required"`
	ReceiverAccount string  `gorm:"size:255;not null;" json:"receiverAccount" binding:"required"`
	Amount          float64 `gorm:"size:255;not null;" json:"amount" binding:"required"`
	TransactionDate string  `gorm:"size:255;not null;" json:"transactionDate" binding:"required"`
	Message         string  `gorm:"size:255;not null;" json:"message"`
}

type CreateTransactionInput struct {
	SenderAccount   string  `gorm:"size:255;not null;" json:"senderAccount" binding:"required"`
	ReceiverAccount string  `gorm:"size:255;not null;" json:"receiverAccount" binding:"required"`
	Amount          float64 `gorm:"size:255;not null;" json:"amount" binding:"required"`
	TransactionDate string  `gorm:"size:255;not null;" json:"transactionDate" binding:"required"`
	Message         string  `gorm:"size:255;not null;" json:"message"`
}
