package account

import (
	"github.com/gin-gonic/gin"
	"github.com/skandansn/webDevBankBackend/entity/bankAccount"
	"github.com/skandansn/webDevBankBackend/models"
	"github.com/skandansn/webDevBankBackend/service/customer"
	"github.com/skandansn/webDevBankBackend/utils"
)

type BankAccountService interface {
	GetBankAccountsByCustomerId(ctx *gin.Context) ([]bankAccount.Account, error)
}

type bankAccountService struct {
}

func (b *bankAccountService) GetBankAccountsByCustomerId(ctx *gin.Context) ([]bankAccount.Account, error) {
	custId := utils.ParseStringAsInt(ctx.GetString("userTypeId"))

	res, err := models.GetBankAccountsByCustomerId(custId)
	if err != nil {
		return []bankAccount.Account{}, err
	}

	return customer.ConvertToAccountDTO(res), nil
}

func NewBankAccountService() BankAccountService {
	return &bankAccountService{}
}
