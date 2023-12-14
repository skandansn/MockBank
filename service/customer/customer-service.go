package customer

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/skandansn/webDevBankBackend/entity/bankAccount"
	"github.com/skandansn/webDevBankBackend/entity/card"
	customerEntity "github.com/skandansn/webDevBankBackend/entity/customer"
	"github.com/skandansn/webDevBankBackend/models"
	"github.com/skandansn/webDevBankBackend/utils"
)

type CustomerService interface {
	GetCustomerDetails(ctx *gin.Context) (customerEntity.Customer, error)
	GetAllCustomers(ctx *gin.Context) ([]customerEntity.Customer, error)
	UpdateProfile(ctx *gin.Context, customer customerEntity.Customer) (customerEntity.Customer, error)
	CreateCustomer(ctx *gin.Context) (customerEntity.Customer, error)
}

type customerService struct {
}

func (c *customerService) CreateCustomer(ctx *gin.Context) (customerEntity.Customer, error) {

	var input customerEntity.CreateCustomerInput

	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		return customerEntity.Customer{}, err
	}

	empId := ctx.MustGet("userTypeId").(string)

	employeeAccess, err := models.IsAccessPresentForEmployee(utils.ParseStringAsInt(empId), "create_customer")

	if err != nil {
		return customerEntity.Customer{}, errors.New("employee does not have access to create customer")
	}

	if !employeeAccess {
		return customerEntity.Customer{}, errors.New("employee does not have access to create customer")
	}

	appointmentId := input.AppointmentId
	customerId := input.CustomerId

	appointment, err := models.GetAppointmentById(appointmentId)
	if err != nil {
		return customerEntity.Customer{}, errors.New("appointment not found")
	}

	if appointment.Status != "Scheduled" || *appointment.CustomerID != customerId {
		return customerEntity.Customer{}, errors.New("appointment is not scheduled for this customer")
	}

	_, err = models.MarkAppointmentAsCompletedByCustIdAndAppointmentId(customerId, appointmentId)

	if err != nil {
		return customerEntity.Customer{}, errors.New("could not mark appointment as completed")
	}

	customer, err := models.GetCustomerById(customerId)

	if err != nil {
		return customerEntity.Customer{}, errors.New("could not find customer")
	}

	for _, account := range input.Accounts {
		dbBankAccount := models.BankAccount{
			CustomerID:     customerId,
			AccountBalance: account.AccountBalance,
			AccountNumber:  utils.GenerateRandomNumberString(10),
			AccountType:    account.AccountType,
		}

		dbBankAccount, err = models.CreateBankAccount(dbBankAccount)

		if dbBankAccount.AccountType == "Checking" {
			dbCard := models.Card{
				CustomerID:        customerId,
				CardNumber:        utils.GenerateRandomNumberString(16),
				CardType:          "Debit",
				CardNetwork:       input.CardNetwork,
				CardCvv:           utils.GenerateRandomNumberString(3),
				CardExpiry:        utils.GenerateCardExpiry(),
				CardName:          "Debit Card",
				CardHolderName:    customer.FirstName + " " + customer.LastName,
				CardLinkedAccount: dbBankAccount.AccountNumber,
			}

			dbCard, err = models.CreateCard(dbCard)

			if err != nil {
				return customerEntity.Customer{}, errors.New("could not create card")
			}
		}

		if err != nil {
			return customerEntity.Customer{}, errors.New("could not create bank account")
		}
	}

	customerDTO := convertToCustomerDTO(customer)

	customerDTO, err = getCardAndAccountDetails(customerDTO)

	if err != nil {
		return customerEntity.Customer{}, errors.New("could not get card and account details")
	}

	return customerDTO, nil
}

func (c *customerService) GetAllCustomers(ctx *gin.Context) ([]customerEntity.Customer, error) {

	empId := ctx.MustGet("userTypeId").(string)

	employeeAccess, err := models.IsAccessPresentForEmployee(utils.ParseStringAsInt(empId), "view_customer_details")

	if err != nil {
		return []customerEntity.Customer{}, errors.New("employee does not have access to view all customers")
	}

	if !employeeAccess {
		return []customerEntity.Customer{}, errors.New("employee does not have access to view all customers")
	}

	dbCustomers, err := models.GetAllCustomers()

	if err != nil {
		return []customerEntity.Customer{}, errors.New("could not get all customers")
	}

	var customers []customerEntity.Customer

	for _, dbCustomer := range dbCustomers {
		customer := convertToCustomerDTO(dbCustomer)
		customer, err = getCardAndAccountDetails(customer)
		if err != nil {
			return []customerEntity.Customer{}, errors.New("could not get card and account details")
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

func (c *customerService) GetCustomerDetails(ctx *gin.Context) (customerEntity.Customer, error) {

	userType := ctx.MustGet("userType").(string)
	userTypeId := ctx.MustGet("userTypeId").(string)

	var dbCustomer models.Customer

	if userType == "employee" {
		employeeAccess, err := models.IsAccessPresentForEmployee(utils.ParseStringAsInt(userTypeId), "view_customer_details")
		if err != nil {
			return customerEntity.Customer{}, errors.New("employee does not have access to view all transactions")
		}
		if !employeeAccess {
			return customerEntity.Customer{}, errors.New("employee does not have access to view all transactions")
		}
	}

	custId := ""
	if userType == "customer" {
		custId = userTypeId
	} else {
		custId = ctx.Query("customerId")
	}

	if custId == "" {
		return customerEntity.Customer{}, errors.New("invalid customer id")
	}

	dbCustomer, err := models.GetCustomerById(utils.ParseStringAsInt(custId))

	if err != nil {
		return customerEntity.Customer{}, errors.New("customer not found")
	}

	customer := convertToCustomerDTO(dbCustomer)

	customer, err = getCardAndAccountDetails(customer)

	if err != nil {
		return customerEntity.Customer{}, errors.New("could not get card and account details")
	}

	return customer, nil
}

func (c *customerService) UpdateProfile(ctx *gin.Context, customer customerEntity.Customer) (customerEntity.Customer, error) {
	userId := ctx.MustGet("userId").(string)

	userIdInt := utils.ParseStringAsInt(userId)

	dbCustomer, err := models.GetCustomerByUserID(userIdInt)

	if err != nil {
		return customerEntity.Customer{}, err
	}

	dbCustomer.FirstName = customer.FirstName
	dbCustomer.LastName = customer.LastName
	dbCustomer.Address = customer.Address
	dbCustomer.DateOfBirth = customer.DateOfBirth
	dbCustomer.Email = customer.Email
	dbCustomer.Phone = customer.Phone
	dbCustomer.UserName = customer.UserName

	dbCustomer, err = models.UpdateCustomer(dbCustomer)

	if err != nil {
		return customerEntity.Customer{}, err
	}

	customer = convertToCustomerDTO(dbCustomer)

	return customer, nil
}

func NewCustomerService() CustomerService {
	return &customerService{}
}

func getCardAndAccountDetails(customer customerEntity.Customer) (customerEntity.Customer, error) {
	customerCards, err := models.GetCardsByCustomerId(customer.CustomerId)

	if err != nil {
		return customerEntity.Customer{}, errors.New("could not get customer cards")
	}

	customerAccounts, err := models.GetBankAccountsByCustomerId(customer.CustomerId)

	if err != nil {
		return customerEntity.Customer{}, errors.New("could not get customer accounts")
	}

	customer.Cards = ConvertToCardDTO(customerCards)
	customer.Accounts = ConvertToAccountDTO(customerAccounts)

	return customer, nil
}

func convertToCustomerDTO(dbCustomer models.Customer) customerEntity.Customer {
	customer := customerEntity.Customer{
		CustomerId:  dbCustomer.ID,
		FirstName:   dbCustomer.FirstName,
		LastName:    dbCustomer.LastName,
		Address:     dbCustomer.Address,
		DateOfBirth: dbCustomer.DateOfBirth,
		Email:       dbCustomer.Email,
		Phone:       dbCustomer.Phone,
		UserName:    dbCustomer.UserName,
	}
	return customer
}

func ConvertToCardDTO(dbCards []models.Card) []card.Card {
	var cards []card.Card
	for _, dbCard := range dbCards {
		card := card.Card{
			CardId:            dbCard.ID,
			CardNumber:        dbCard.CardNumber,
			CardType:          dbCard.CardType,
			CardNetwork:       dbCard.CardNetwork,
			CardCvv:           dbCard.CardCvv,
			CardExpiry:        dbCard.CardExpiry,
			CardName:          dbCard.CardName,
			CardHolderName:    dbCard.CardHolderName,
			CardLinkedAccount: dbCard.CardLinkedAccount,
			CardLimit: 		   dbCard.CardLimit,
			CustomerId:        dbCard.CustomerID,
		}
		cards = append(cards, card)
	}
	return cards
}

func ConvertToAccountDTO(dbAccounts []models.BankAccount) []bankAccount.Account {
	var accounts []bankAccount.Account
	for _, dbAccount := range dbAccounts {
		account := bankAccount.Account{
			AccountID:      dbAccount.ID,
			AccountBalance: dbAccount.AccountBalance,
			AccountNumber:  dbAccount.AccountNumber,
			AccountType:    dbAccount.AccountType,
			CustomerId:     dbAccount.CustomerID,
		}
		accounts = append(accounts, account)
	}
	return accounts
}
