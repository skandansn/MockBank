package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/skandansn/webDevBankBackend/entity/appointment"
	"github.com/skandansn/webDevBankBackend/entity/auth"
	customerEntity "github.com/skandansn/webDevBankBackend/entity/customer"
	"github.com/skandansn/webDevBankBackend/models"
	"github.com/skandansn/webDevBankBackend/utils/token"
	"net/http"
)

func Login(c *gin.Context) {

	var input auth.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password

	token, err := models.LoginCheck(u.Username, u.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func Register(c *gin.Context) (customerEntity.Customer, error) {

	var input auth.CustomerRegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		return customerEntity.Customer{}, err
	}

	cus, err := coreCustomerRegister(input)

	if err != nil {
		return customerEntity.Customer{}, err
	}

	cusReturn := convertToCustomerDTO(cus)

	return cusReturn, nil
}

func RegisterFromAppointment(bookAppointmentInput appointment.BookAppointmentInput) (customerEntity.Customer, error) {

	var input auth.CustomerRegisterInput

	input.FirstName = bookAppointmentInput.CustomerDetails.FirstName
	input.LastName = bookAppointmentInput.CustomerDetails.LastName
	input.Address = bookAppointmentInput.CustomerDetails.Address
	input.DateOfBirth = bookAppointmentInput.CustomerDetails.DateOfBirth
	input.Email = bookAppointmentInput.CustomerDetails.Email
	input.Phone = bookAppointmentInput.CustomerDetails.Phone
	input.UserName = bookAppointmentInput.CustomerDetails.UserName
	input.Password = bookAppointmentInput.CustomerDetails.Password

	cus, err := coreCustomerRegister(input)

	if err != nil {
		return customerEntity.Customer{}, err
	}

	return convertToCustomerDTO(cus), nil
}

func CurrentUser(c *gin.Context) {

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

func coreCustomerRegister(input auth.CustomerRegisterInput) (models.Customer, error) {

	u := models.User{}
	cus := models.Customer{}

	u.Username = input.UserName
	u.Password = input.Password
	u.Type = "customer"

	cus.FirstName = input.FirstName
	cus.LastName = input.LastName
	cus.Address = input.Address
	cus.DateOfBirth = input.DateOfBirth
	cus.Email = input.Email
	cus.Phone = input.Phone
	cus.UserName = input.UserName

	_, err := u.SaveUser()

	if err != nil {
		return models.Customer{}, err
	}

	cus.UserID = u.ID

	_, err = cus.SaveCustomer()

	if err != nil {
		return models.Customer{}, err
	}

	u.UserTypeID = cus.ID

	_, err = u.UpdateUser()

	if err != nil {
		return models.Customer{}, err
	}

	return cus, nil
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
