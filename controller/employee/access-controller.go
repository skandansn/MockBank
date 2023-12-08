package employee

import (
	"github.com/gin-gonic/gin"
	"github.com/skandansn/webDevBankBackend/entity/employee"
	employeeService "github.com/skandansn/webDevBankBackend/service/employee"
)

type AccessController interface {
	GetAccessForEmployee(employeeID string) ([]employee.Access, error)
	SaveAccessForEmployee(employeeID string, ctx *gin.Context) ([]employee.Access, error)
}

type accessController struct {
	service employeeService.AccessService
}

func (c *accessController) GetAccessForEmployee(employeeID string) ([]employee.Access, error) {
	res, err := c.service.GetAccessForEmployee(employeeID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *accessController) SaveAccessForEmployee(employeeID string, ctx *gin.Context) ([]employee.Access, error) {
	var accessList employee.AccessUpdateInput
	ctx.BindJSON(&accessList)

	res, err := c.service.SaveAccessForEmployee(employeeID, accessList.AccessList)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewAccessController(service employeeService.AccessService) AccessController {
	return &accessController{
		service: service,
	}
}
