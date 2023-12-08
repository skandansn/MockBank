package customer

import (
	"github.com/gin-gonic/gin"
	customerEntity "github.com/skandansn/webDevBankBackend/entity/customer"
	"github.com/skandansn/webDevBankBackend/service/customer"
)

type CustomerController interface {
	GetCustomerDetails(ctx *gin.Context) (customerEntity.Customer, error)
	UpdateProfile(ctx *gin.Context) (customerEntity.Customer, error)
	CreateCustomer(ctx *gin.Context) (customerEntity.Customer, error)
}

type customerController struct {
	service customer.CustomerService
}

func (c *customerController) CreateCustomer(ctx *gin.Context) (customerEntity.Customer, error) {
	res, err := c.service.CreateCustomer(ctx)
	if err != nil {
		return customerEntity.Customer{}, err
	}
	return res, nil
}

func (c *customerController) GetCustomerDetails(ctx *gin.Context) (customerEntity.Customer, error) {
	res, err := c.service.GetCustomerDetails(ctx)
	if err != nil {
		return customerEntity.Customer{}, err
	}
	return res, nil
}

func (c *customerController) UpdateProfile(ctx *gin.Context) (customerEntity.Customer, error) {
	var input customerEntity.Customer
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		return customerEntity.Customer{}, err
	}

	res, err := c.service.UpdateProfile(ctx, input)
	if err != nil {
		return customerEntity.Customer{}, err
	}
	return res, nil
}

func NewCustomerController(service customer.CustomerService) CustomerController {
	return &customerController{
		service: service,
	}
}
