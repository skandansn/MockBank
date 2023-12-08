package card

import (
	"github.com/gin-gonic/gin"
	"github.com/skandansn/webDevBankBackend/entity/card"
	cardService "github.com/skandansn/webDevBankBackend/service/card"
)

type CardController interface {
	GetCardsByCustomerId(ctx *gin.Context) ([]card.Card, error)
	CreateCardRequest(ctx *gin.Context) (card.CardRequest, error)
	GetPendingCardRequests(ctx *gin.Context) ([]card.CardRequest, error)
	ApproveOrRejectCardRequest(ctx *gin.Context) (card.CardRequest, error)
}

type cardController struct {
	service cardService.CardService
}

func (c *cardController) ApproveOrRejectCardRequest(ctx *gin.Context) (card.CardRequest, error) {
	res, err := c.service.ApproveOrRejectCardRequest(ctx)
	if err != nil {
		return card.CardRequest{}, err
	}

	return res, nil
}

func (c *cardController) GetPendingCardRequests(ctx *gin.Context) ([]card.CardRequest, error) {
	res, err := c.service.GetPendingCardRequests(ctx)
	if err != nil {
		return []card.CardRequest{}, err
	}

	return res, nil
}

func (c *cardController) CreateCardRequest(ctx *gin.Context) (card.CardRequest, error) {
	res, err := c.service.CreateCardRequest(ctx)
	if err != nil {
		return card.CardRequest{}, err
	}

	return res, nil
}

func (c *cardController) GetCardsByCustomerId(ctx *gin.Context) ([]card.Card, error) {
	res, err := c.service.GetCardsByCustomerId(ctx)
	if err != nil {
		return []card.Card{}, err
	}
	return res, nil
}

func NewCardController(service cardService.CardService) CardController {
	return &cardController{
		service: service,
	}
}
