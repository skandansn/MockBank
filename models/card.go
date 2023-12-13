package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type Card struct {
	gorm.Model
	CustomerID        uint   `gorm:"not null" json:"customer_id"`
	CardNumber        string `gorm:"size:255;not null;" json:"card_number"`
	CardType          string `gorm:"size:255;not null;" json:"card_type"`
	CardNetwork       string `gorm:"size:255;not null;" json:"card_network"`
	CardCvv           string `gorm:"size:255;not null;" json:"card_cvv"`
	CardExpiry        string `gorm:"size:255;not null;" json:"card_expiry"`
	CardName          string `gorm:"size:255;not null;" json:"card_name"`
	CardHolderName    string `gorm:"size:255;not null;" json:"card_holder_name"`
	CardLimit         uint   `gorm:"" json:"card_limit"`
	CardLinkedAccount string `gorm:"size:255;" json:"card_linked_account"`
}

type CardRequest struct {
	gorm.Model
	CustomerID  uint   `gorm:"not null" json:"customer_id"`
	CardType    string `gorm:"size:255;not null;" json:"card_type"`
	CardNetwork string `gorm:"size:255;not null;" json:"card_network"`
	CardLimit   uint   `gorm:"" json:"card_limit"`
	CardName    string `gorm:"size:255;not null;" json:"card_name"`
	Status      string `gorm:"size:255;not null;" json:"status"`
	Reason      string `gorm:"size:255;" json:"reason"`
}

func CreateCard(c Card) (Card, error) {

	err := DB.Create(&c).Error
	if err != nil {
		return Card{}, err
	}

	return c, nil
}

func GetCardsByCustomerId(cid uint) ([]Card, error) {

	var c []Card

	if err := DB.Where("customer_id = ?", cid).Find(&c).Error; err != nil {
		return c, err
	}

	return c, nil
}

func CreateCardRequest(c CardRequest) (CardRequest, error) {

	err := DB.Create(&c).Error
	if err != nil {
		return CardRequest{}, err
	}

	return c, nil
}

func GetPendingCardRequests() ([]CardRequest, error) {

	var c []CardRequest

	if err := DB.Where("status = ?", "pending").Find(&c).Error; err != nil {
		return c, err
	}

	return c, nil
}

func GetPendingCardRequestsForCustomer(cid uint) ([]CardRequest, error) {

	var c []CardRequest

	if err := DB.Where("status = ? AND customer_id = ?", "pending", cid).Find(&c).Error; err != nil {
		return c, err
	}

	return c, nil
}

func UpdateCardRequestStatus(cid uint, status string, reason string) (CardRequest, error) {

	var c CardRequest

	if err := DB.Where("id = ?", cid).First(&c).Error; err != nil {
		return c, err
	}

	if c.Status != "pending" {
		return c, errors.New("card request is not pending")
	}

	c.Status = status
	c.Reason = reason

	err := DB.Save(&c).Error
	if err != nil {
		return CardRequest{}, err
	}

	return c, nil
}
