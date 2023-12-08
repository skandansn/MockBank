package employee

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	employeeEntity "github.com/skandansn/webDevBankBackend/entity/employee"
	employeeService "github.com/skandansn/webDevBankBackend/service/employee"
	bankValidators "github.com/skandansn/webDevBankBackend/validators"
)

type EmployeeController interface {
	FindAll() ([]employeeEntity.Employee, error)
	Save(ctx *gin.Context) (employeeEntity.Employee, error)
	Edit(ctx *gin.Context) (employeeEntity.Employee, error)
	Delete(employeeId string) error
}

type controller struct {
	service employeeService.EmployeeService
}

func (c *controller) FindAll() ([]employeeEntity.Employee, error) {
	res, err := c.service.FindAll()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *controller) Save(ctx *gin.Context) (employeeEntity.Employee, error) {
	var employee employeeEntity.EmployeeRegisterInput
	err := ctx.BindJSON(&employee)
	if err != nil {
		return employeeEntity.Employee{}, err
	}
	err = validate.Struct(employee)
	if err != nil {
		return employeeEntity.Employee{}, err
	}
	res, err := c.service.Save(employee)
	if err != nil {
		return employeeEntity.Employee{}, err
	}
	return res, nil
}

func (c *controller) Edit(ctx *gin.Context) (employeeEntity.Employee, error) {
	var employee employeeEntity.Employee
	ctx.BindJSON(&employee)
	res, err := c.service.Edit(employee)
	if err != nil {
		return employeeEntity.Employee{}, err
	}
	return res, nil
}

func (c *controller) Delete(employeeId string) error {
	err := c.service.Delete(employeeId)
	if err != nil {
		return err
	}
	return nil
}

var validate *validator.Validate

func New(service employeeService.EmployeeService) EmployeeController {
	validate = validator.New()
	validate.RegisterValidation("isValidRole", bankValidators.ValidateEmployeeRole)
	return &controller{
		service: service,
	}
}
