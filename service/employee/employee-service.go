package employee

import (
	"fmt"
	employeeEntity "github.com/skandansn/webDevBankBackend/entity/employee"
	"github.com/skandansn/webDevBankBackend/models"
	"strconv"
)

type EmployeeService interface {
	FindAll() ([]employeeEntity.Employee, error)
	Save(employee employeeEntity.EmployeeRegisterInput) (employeeEntity.Employee, error)
	Edit(employee employeeEntity.Employee) (employeeEntity.Employee, error)
	Delete(employeeID string) error
}

type employeeService struct {
}

func (service *employeeService) FindAll() ([]employeeEntity.Employee, error) {
	var employees []employeeEntity.Employee

	dbEmployees, err := models.GetAllEmployees()

	if err != nil {
		return nil, err
	}

	for _, v := range dbEmployees {
		currentEmployee := convertToEmployeeDTO(v)
		accesses, err := models.GetAllAccessForEmployee(v.ID)
		if err != nil {
			return nil, err
		}
		currentEmployee.AccessList = convertToAccessDTO(accesses)
		employees = append(employees, currentEmployee)
	}

	return employees, nil
}

func (service *employeeService) Save(employee employeeEntity.EmployeeRegisterInput) (employeeEntity.Employee, error) {

	user := models.User{
		Username: employee.Username,
		Password: employee.Password,
		Type:     "employee",
	}

	_, err := user.SaveUser()

	if err != nil {
		return employeeEntity.Employee{}, err
	}

	emp := models.Employee{
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		Branch:    employee.Branch,
		Role:      employee.Role,
		Phone:     employee.Phone,
		Email:     employee.Email,
		UserName:  employee.Username,
		UserID:    user.ID,
	}

	dbEmployee, err := models.SaveEmployee(emp)

	if err != nil {
		return employeeEntity.Employee{}, err
	}

	user.UserTypeID = dbEmployee.ID

	_, err = user.UpdateUser()

	if err != nil {
		return employeeEntity.Employee{}, err
	}

	employeeAccesses := employee.AccessList
	employeeAccessItems := employeeEntity.GetEmployeeAccessItems()

	accessModels := make([]models.Access, 0)

	for _, v := range employeeAccesses {
		if employeeAccessItems[v.AccessName] {
			accessModels = append(accessModels, models.Access{
				EmployeeID:  dbEmployee.ID,
				AccessName:  v.AccessName,
				AccessGiven: v.AccessGiven,
			})
		}
	}

	dbAccess := []models.Access{}

	if len(accessModels) != 0 {
		dbAccess, err = models.SaveAccessesForEmployee(accessModels)
		if err != nil {
			return employeeEntity.Employee{}, err
		}
	}

	empReturn := convertToEmployeeDTO(*dbEmployee)

	empReturn.AccessList = convertToAccessDTO(dbAccess)

	return empReturn, nil
}

func (service *employeeService) Edit(employee employeeEntity.Employee) (employeeEntity.Employee, error) {
	employeeIDUint := employee.EmployeeID

	dbEmployee, err := models.GetEmployeeById(employeeIDUint)
	if err != nil {
		return employeeEntity.Employee{}, err
	}

	replaceNonEmptyEmployeeFields(&dbEmployee, employee)

	_, err = models.UpdateEmployee(dbEmployee)

	if err != nil {
		return employeeEntity.Employee{}, err
	}

	accesses, err := models.GetAllAccessForEmployee(employeeIDUint)
	if err != nil {
		return employeeEntity.Employee{}, err
	}
	empReturn := convertToEmployeeDTO(dbEmployee)
	empReturn.AccessList = convertToAccessDTO(accesses)
	return empReturn, nil
}

func (service *employeeService) Delete(employeeID string) error {
	empIDUint, err := strconv.ParseUint(employeeID, 10, 0)
	if err != nil {
		fmt.Println("Error converting string to uint:", err)
		return err
	}
	err = models.DeleteEmployee(uint(empIDUint))
	if err != nil {
		return err
	}
	return nil
}

func replaceNonEmptyEmployeeFields(employee *models.Employee, newEmployee employeeEntity.Employee) {
	if newEmployee.FirstName != "" {
		employee.FirstName = newEmployee.FirstName
	}
	if newEmployee.LastName != "" {
		employee.LastName = newEmployee.LastName
	}
	if newEmployee.Branch != "" {
		employee.Branch = newEmployee.Branch
	}
	if newEmployee.Role != "" {
		employee.Role = newEmployee.Role
	}
	if newEmployee.Phone != "" {
		employee.Phone = newEmployee.Phone
	}
	if newEmployee.Email != "" {
		employee.Email = newEmployee.Email
	}
}

func New() EmployeeService {
	return &employeeService{}
}

func convertToEmployeeDTO(emp models.Employee) employeeEntity.Employee {
	return employeeEntity.Employee{
		FirstName:  emp.FirstName,
		LastName:   emp.LastName,
		Branch:     emp.Branch,
		Role:       emp.Role,
		EmployeeID: emp.ID,
		Phone:      emp.Phone,
		Email:      emp.Email,
	}
}

func convertToAccessDTO(accesses []models.Access) []employeeEntity.Access {
	accessList := make([]employeeEntity.Access, 0)
	for _, v := range accesses {
		accessList = append(accessList, employeeEntity.Access{
			AccessName:  v.AccessName,
			AccessGiven: v.AccessGiven,
		})
	}
	return accessList
}
