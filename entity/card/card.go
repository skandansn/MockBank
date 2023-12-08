package card

type Card struct {
	CardId            uint   `gorm:"primary_key;auto_increment" json:"cardId" binding:"required"`
	CustomerId        uint   `gorm:"not null" json:"customerId" binding:"required"`
	CardNumber        string `gorm:"size:255;not null;" json:"cardNumber" binding:"required"`
	CardType          string `gorm:"size:255;not null;" json:"cardType" binding:"required"`
	CardNetwork       string `gorm:"size:255;not null;" json:"cardNetwork" binding:"required"`
	CardCvv           string `gorm:"size:255;not null;" json:"cardCvv" binding:"required"`
	CardExpiry        string `gorm:"size:255;not null;" json:"cardExpiry" binding:"required"`
	CardName          string `gorm:"size:255;not null;" json:"cardName" binding:"required"`
	CardHolderName    string `gorm:"size:255;not null;" json:"cardHolderName" binding:"required"`
	CardLimit         uint   `gorm:"size:255;not null;" json:"cardLimit" binding:"required"`
	CardLinkedAccount string `gorm:"size:255;not null;" json:"cardLinkedAccount" binding:"required"`
}

type CardRequestInput struct {
	CardType    string `gorm:"size:255;not null;" json:"cardType" binding:"required"`
	CardNetwork string `gorm:"size:255;not null;" json:"cardNetwork" binding:"required"`
	CardLimit   uint   `gorm:"size:255;not null;" json:"cardLimit" binding:"required"`
	CardName    string `gorm:"size:255;not null;" json:"cardName" binding:"required"`
}

type CardRequest struct {
	CardRequestId uint   `gorm:"primary_key;auto_increment" json:"cardRequestId" binding:"required"`
	CustomerId    uint   `gorm:"not null" json:"customerId"`
	CardType      string `gorm:"size:255;not null;" json:"cardType" binding:"required"`
	CardNetwork   string `gorm:"size:255;not null;" json:"cardNetwork" binding:"required"`
	CardLimit     uint   `gorm:"size:255;not null;" json:"cardLimit" binding:"required"`
	CardName      string `gorm:"size:255;not null;" json:"cardName" binding:"required"`
	Status        string `gorm:"size:255;not null;" json:"status"`
	Reason        string `gorm:"size:255;" json:"reason"`
}

type CardRequestDecision struct {
	CardRequestId uint   `gorm:"primary_key;auto_increment" json:"cardRequestId" binding:"required"`
	Decision      string `gorm:"size:255;not null;" json:"decision" binding:"required"`
	Reason        string `gorm:"size:255;" json:"reason"`
}
