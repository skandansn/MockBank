package transaction

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/skandansn/webDevBankBackend/entity/transactions"
	"github.com/skandansn/webDevBankBackend/models"
	"github.com/skandansn/webDevBankBackend/utils"
)

type TransactionService interface {
	CreateTransaction(ctx *gin.Context) (transactions.Transaction, error)
	GetTransactionsForSenderAccount(ctx *gin.Context) ([]transactions.Transaction, error)
	GetTransactionsForCustomer(ctx *gin.Context) ([]transactions.Transaction, error)
	GetTransactionsForCustomerByEmployee(ctx *gin.Context) ([]transactions.Transaction, error)
}

type transactionService struct {
}

func (t *transactionService) GetTransactionsForCustomerByEmployee(ctx *gin.Context) ([]transactions.Transaction, error) {
	empId := ctx.GetString("userTypeId")
	if empId == "" {
		return []transactions.Transaction{}, errors.New("invalid employee id")
	}

	employeeAccessAvailable, err := models.IsAccessPresentForEmployee(utils.ParseStringAsInt(empId), "view_customer_transactions")
	if err != nil {
		return []transactions.Transaction{}, errors.New("employee does not have access to view all transactions")
	}

	if !employeeAccessAvailable {
		return []transactions.Transaction{}, errors.New("employee does not have access to view all transactions")
	}

	custId := ctx.Param("customerId")
	if custId == "" {
		return []transactions.Transaction{}, errors.New("invalid customer id")
	}

	customerAccounts, err := models.GetBankAccountsByCustomerId(utils.ParseStringAsInt(custId))

	if err != nil {
		return []transactions.Transaction{}, err
	}

	var customerAccountNumbers []string

	for _, v := range customerAccounts {
		customerAccountNumbers = append(customerAccountNumbers, v.AccountNumber)
	}

	customerTransactions, err := models.GetTransactionsForBankAccounts(customerAccountNumbers, "1970-01-01", utils.GetCurrentDate())

	if err != nil {
		return []transactions.Transaction{}, err
	}

	return convertToTransactionDTOs(customerTransactions), nil

}

func (t *transactionService) GetTransactionsForCustomer(ctx *gin.Context) ([]transactions.Transaction, error) {
	custId := ctx.GetString("userTypeId")
	if custId == "" {
		return []transactions.Transaction{}, errors.New("invalid customer id")
	}

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	if endDate == "" {
		endDate = utils.GetCurrentDate()
	}

	if startDate == "" {
		startDate = "1970-01-01"
	}

	customerBankAccounts, err := models.GetBankAccountsByCustomerId(utils.ParseStringAsInt(custId))

	if err != nil {
		return []transactions.Transaction{}, err
	}

	var customerBankAccountNumbers []string

	for _, v := range customerBankAccounts {
		customerBankAccountNumbers = append(customerBankAccountNumbers, v.AccountNumber)
	}

	var res []models.Transaction

	res, err = models.GetTransactionsForBankAccounts(customerBankAccountNumbers, startDate, endDate)

	if err != nil {
		return []transactions.Transaction{}, err
	}

	return convertToTransactionDTOs(res), nil
}

func (t *transactionService) GetTransactionsForSenderAccount(ctx *gin.Context) ([]transactions.Transaction, error) {
	custId := ctx.GetString("userTypeId")
	if custId == "" {
		return []transactions.Transaction{}, errors.New("invalid customer id")
	}

	accountId := ctx.Query("accountId")

	if accountId == "" {
		return []transactions.Transaction{}, errors.New("invalid account id")
	}

	count := ctx.Query("count")

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	if endDate == "" {
		endDate = utils.GetCurrentDate()
	}

	if startDate == "" {
		startDate = "1970-01-01"
	}

	_, err := models.CheckIfAccountBelongsToCustomer(accountId, utils.ParseStringAsInt(custId))

	if err != nil {
		return []transactions.Transaction{}, err
	}

	var res []models.Transaction

	if count != "" {
		res, err = models.GetNLatestTransactionsForSenderAccount(accountId, utils.ParseStringAsInt(count), startDate, endDate)
	} else {
		res, err = models.GetTransactionsForSenderAccount(accountId, startDate, endDate)
	}

	if err != nil {
		return []transactions.Transaction{}, err
	}

	return convertToTransactionDTOs(res), nil
}

func (t *transactionService) CreateTransaction(ctx *gin.Context) (transactions.Transaction, error) {
	var input transactions.CreateTransactionInput
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		return transactions.Transaction{}, err
	}

	custId := ctx.GetString("userTypeId")
	if custId == "" {
		return transactions.Transaction{}, errors.New("invalid customer id")
	}

	if input.SenderAccount == input.ReceiverAccount {
		return transactions.Transaction{}, errors.New("sender and receiver accounts cannot be same")
	}

	_, err = models.CheckIfAccountBelongsToCustomerAndBalance(input.SenderAccount, utils.ParseStringAsInt(custId), input.Amount)

	if err != nil {
		return transactions.Transaction{}, err
	}

	_, err = models.CheckIfAccountExists(input.ReceiverAccount)

	if err != nil {
		return transactions.Transaction{}, err
	}

	res, err := models.CreateTransaction(connvertToTransactionDB(input))

	if err != nil {
		return transactions.Transaction{}, err
	}

	_, err = models.UpdateSenderAndReceiverAccountBalance(input.SenderAccount, input.ReceiverAccount, input.Amount)

	if err != nil {
		return transactions.Transaction{}, err
	}

	return convertToTransactionDTO(res), nil
}

func connvertToTransactionDB(input transactions.CreateTransactionInput) models.Transaction {
	return models.Transaction{
		SenderAccount:   input.SenderAccount,
		ReceiverAccount: input.ReceiverAccount,
		Amount:          input.Amount,
		TransactionDate: input.TransactionDate,
		Message:         input.Message,
	}
}

func convertToTransactionDTO(input models.Transaction) transactions.Transaction {
	return transactions.Transaction{
		TransactionId:   input.ID,
		SenderAccount:   input.SenderAccount,
		ReceiverAccount: input.ReceiverAccount,
		Amount:          input.Amount,
		TransactionDate: input.TransactionDate,
		Message:         input.Message,
	}
}

func convertToTransactionDTOs(input []models.Transaction) []transactions.Transaction {
	var res []transactions.Transaction
	for _, v := range input {
		res = append(res, convertToTransactionDTO(v))
	}
	return res
}

func NewTransactionService() TransactionService {
	return &transactionService{}
}
