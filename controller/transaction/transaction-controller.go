package transaction

import (
	"github.com/gin-gonic/gin"
	"github.com/skandansn/webDevBankBackend/entity/transactions"
	transactionService "github.com/skandansn/webDevBankBackend/service/transaction"
)

type TransactionController interface {
	CreateTransaction(ctx *gin.Context) (transactions.Transaction, error)
	GetTransactionsForSenderAccount(ctx *gin.Context) ([]transactions.Transaction, error)
	GetTransactionsForCustomer(ctx *gin.Context) ([]transactions.Transaction, error)
	GetTransactionsForCustomerByEmployee(ctx *gin.Context) ([]transactions.Transaction, error)
}

type transactionController struct {
	service transactionService.TransactionService
}

func (t *transactionController) GetTransactionsForCustomerByEmployee(ctx *gin.Context) ([]transactions.Transaction, error) {
	res, err := t.service.GetTransactionsForCustomerByEmployee(ctx)
	if err != nil {
		return []transactions.Transaction{}, err
	}
	return res, nil
}

func (t *transactionController) GetTransactionsForCustomer(ctx *gin.Context) ([]transactions.Transaction, error) {
	res, err := t.service.GetTransactionsForCustomer(ctx)
	if err != nil {
		return []transactions.Transaction{}, err
	}
	return res, nil
}

func (t *transactionController) GetTransactionsForSenderAccount(ctx *gin.Context) ([]transactions.Transaction, error) {
	res, err := t.service.GetTransactionsForSenderAccount(ctx)
	if err != nil {
		return []transactions.Transaction{}, err
	}
	return res, nil
}

func (t *transactionController) CreateTransaction(ctx *gin.Context) (transactions.Transaction, error) {
	res, err := t.service.CreateTransaction(ctx)
	if err != nil {
		return transactions.Transaction{}, err
	}

	return res, nil
}

func NewTransactionController(service transactionService.TransactionService) TransactionController {
	return &transactionController{
		service: service,
	}
}
