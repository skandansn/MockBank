package account

import (
	"github.com/gin-gonic/gin"
	"github.com/skandansn/webDevBankBackend/entity/bankAccount"
	"github.com/skandansn/webDevBankBackend/service/account"
)

type BankAccountController interface {
	GetBankAccountsByCustomerId(ctx *gin.Context) ([]bankAccount.Account, error)
}

type bankAccountController struct {
	service account.BankAccountService
}

func (b *bankAccountController) GetBankAccountsByCustomerId(ctx *gin.Context) ([]bankAccount.Account, error) {
	res, err := b.service.GetBankAccountsByCustomerId(ctx)
	if err != nil {
		return []bankAccount.Account{}, err
	}
	return res, nil
}

func NewBankAccountController(service account.BankAccountService) BankAccountController {
	return &bankAccountController{
		service: service,
	}
}
