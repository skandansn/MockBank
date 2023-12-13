package card

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/skandansn/webDevBankBackend/entity/card"
	"github.com/skandansn/webDevBankBackend/models"
	"github.com/skandansn/webDevBankBackend/service/customer"
	"github.com/skandansn/webDevBankBackend/utils"
)

type CardService interface {
	GetCardsByCustomerId(ctx *gin.Context) ([]card.Card, error)
	CreateCardRequest(ctx *gin.Context) (card.CardRequest, error)
	GetPendingCardRequests(ctx *gin.Context) ([]card.CardRequest, error)
	ApproveOrRejectCardRequest(ctx *gin.Context) (card.CardRequest, error)
}

type cardService struct {
}

func (c cardService) ApproveOrRejectCardRequest(ctx *gin.Context) (card.CardRequest, error) {
	err := cardRequestPermissionCheck(ctx)
	if err != nil {
		return card.CardRequest{}, err
	}

	var cardReqDecision card.CardRequestDecision

	err = ctx.ShouldBindJSON(&cardReqDecision)
	if err != nil {
		return card.CardRequest{}, err
	}

	if cardReqDecision.Decision != "approved" && cardReqDecision.Decision != "rejected" {
		return card.CardRequest{}, errors.New("invalid decision")
	}

	cardReqDb, err := models.UpdateCardRequestStatus(cardReqDecision.CardRequestId, cardReqDecision.Decision, cardReqDecision.Reason)

	if err != nil {
		return card.CardRequest{}, err
	}

	if cardReqDecision.Decision == "approved" {
		cardDb := models.Card{
			CustomerID:        cardReqDb.CustomerID,
			CardNumber:        utils.GenerateRandomNumberString(16),
			CardType:          cardReqDb.CardType,
			CardNetwork:       cardReqDb.CardNetwork,
			CardCvv:           utils.GenerateRandomNumberString(3),
			CardExpiry:        utils.GenerateCardExpiry(),
			CardName:          cardReqDb.CardName,
			CardHolderName:    cardReqDb.CardName,
			CardLimit:         cardReqDb.CardLimit,
			CardLinkedAccount: "",
		}

		_, err = models.CreateCard(cardDb)
		if err != nil {
			return card.CardRequest{}, errors.New("error creating card")
		}
	}

	return convertToCardRequestDTO(cardReqDb), nil
}

func (c cardService) GetPendingCardRequests(ctx *gin.Context) ([]card.CardRequest, error) {

	isUserCustomer := ctx.GetString("userType") == "customer"

	if isUserCustomer {
		custId := utils.ParseStringAsInt(ctx.GetString("userTypeId"))

		res, err := models.GetPendingCardRequestsForCustomer(custId)

		if err != nil {
			return []card.CardRequest{}, errors.New("error getting pending card requests")
		}

		cardRequests := []card.CardRequest{}
		for _, cardReq := range res {
			cardRequests = append(cardRequests, convertToCardRequestDTO(cardReq))
		}

		return cardRequests, nil
	}

	err := cardRequestPermissionCheck(ctx)
	if err != nil {
		return []card.CardRequest{}, err
	}

	res, err := models.GetPendingCardRequests()
	if err != nil {
		return []card.CardRequest{}, errors.New("error getting pending card requests")
	}

	cardRequests := []card.CardRequest{}
	for _, cardReq := range res {
		cardRequests = append(cardRequests, convertToCardRequestDTO(cardReq))
	}

	return cardRequests, nil
}

func (c cardService) GetCardsByCustomerId(ctx *gin.Context) ([]card.Card, error) {
	custId := utils.ParseStringAsInt(ctx.GetString("userTypeId"))

	res, err := models.GetCardsByCustomerId(custId)
	if err != nil {
		return []card.Card{}, err
	}

	return customer.ConvertToCardDTO(res), nil
}

func (c cardService) CreateCardRequest(ctx *gin.Context) (card.CardRequest, error) {
	var cardReq card.CardRequestInput

	err := ctx.ShouldBindJSON(&cardReq)
	if err != nil {
		return card.CardRequest{}, err
	}

	custId := ctx.GetString("userTypeId")

	if custId == "" {
		return card.CardRequest{}, errors.New("invalid customer id")
	}

	custIdInt := utils.ParseStringAsInt(custId)

	cardReqDb := convertToCardRequestInputDB(cardReq)

	cardReqDb.CustomerID = uint(custIdInt)

	res, err := models.CreateCardRequest(cardReqDb)

	if err != nil {
		return card.CardRequest{}, errors.New("error creating card request")
	}

	return convertToCardRequestDTO(res), nil
}

func convertToCardRequestInputDB(c card.CardRequestInput) models.CardRequest {
	return models.CardRequest{
		CardType:    c.CardType,
		CardNetwork: c.CardNetwork,
		CardLimit:   c.CardLimit,
		CardName:    c.CardName,
		Status:      "pending",
		Reason:      "",
	}
}

func convertToCardRequestDTO(c models.CardRequest) card.CardRequest {
	return card.CardRequest{
		CardRequestId: c.ID,
		CustomerId:    c.CustomerID,
		CardType:      c.CardType,
		CardNetwork:   c.CardNetwork,
		CardLimit:     c.CardLimit,
		CardName:      c.CardName,
		Status:        c.Status,
		Reason:        c.Reason,
	}
}

func cardRequestPermissionCheck(ctx *gin.Context) error {
	employeeId := utils.ParseStringAsInt(ctx.GetString("userTypeId"))

	employeeAccesses, err := models.IsAccessPresentForEmployee(employeeId, "approve_card")

	if err != nil {
		return errors.New("error getting employee accesses or employee does not have access to approve card requests")
	}

	if !employeeAccesses {
		return errors.New("employee does not have access to approve card requests")
	}

	return nil
}

func NewCardService() CardService {
	return &cardService{}
}
