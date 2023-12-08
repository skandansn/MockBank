package employee

import (
	"errors"
	"github.com/skandansn/webDevBankBackend/entity/employee"
	"github.com/skandansn/webDevBankBackend/models"
	"github.com/skandansn/webDevBankBackend/utils"
)

type AccessService interface {
	GetAccessForEmployee(employeeID string) ([]employee.Access, error)
	SaveAccessForEmployee(employeeID string, accessList []employee.Access) ([]employee.Access, error)
}

type accessService struct {
}

func (service *accessService) GetAccessForEmployee(employeeID string) ([]employee.Access, error) {
	employeeIDUint := utils.ParseStringAsInt(employeeID)

	var accessList []employee.Access

	dbAccess, err := models.GetAllAccessForEmployee(employeeIDUint)

	if err != nil {
		return nil, err
	}

	accessList = convertToAccessDTO(dbAccess)

	return accessList, nil
}

func (service *accessService) SaveAccessForEmployee(employeeID string, accessList []employee.Access) ([]employee.Access, error) {
	employeeIDUint := utils.ParseStringAsInt(employeeID)

	employeeAccessItems := employee.GetEmployeeAccessItems()

	accessModels := make([]models.Access, 0)

	for _, v := range accessList {
		if employeeAccessItems[v.AccessName] {
			accessModels = append(accessModels, models.Access{
				EmployeeID:  employeeIDUint,
				AccessName:  v.AccessName,
				AccessGiven: v.AccessGiven,
			})
		} else {
			return []employee.Access{}, errors.New("Invalid access name" + v.AccessName + " for employee")
		}
	}

	dbAccess := []models.Access{}

	if len(accessModels) != 0 {
		_, err := models.SaveAccessesForEmployee(accessModels)
		if err != nil {
			return []employee.Access{}, err
		}
	}

	return convertToAccessDTO(dbAccess), nil
}

func NewAccessService() AccessService {
	return &accessService{}
}
