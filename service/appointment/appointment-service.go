package appointment

import (
	"encoding/json"
	"errors"
	"github.com/skandansn/webDevBankBackend/entity/appointment"
	customerEntity "github.com/skandansn/webDevBankBackend/entity/customer"
	"github.com/skandansn/webDevBankBackend/models"
	"github.com/skandansn/webDevBankBackend/utils"
)

type AppointmentService interface {
	GetAvailableAppointments(date string) ([]appointment.AvailableAppointment, error)
	CreateAppointment(appointment appointment.Appointment) (appointment.Appointment, error)
	GetAppointmentById(id uint) (appointment.Appointment, error)
	GetAppointmentByCustIdAndId(custId interface{}, id uint) (appointment.Appointment, error)
	GetBookedAppointmentsForCustomer(custId interface{}) ([]appointment.Appointment, error)
	GetAppointmentsForEmployee(employeeId uint) ([]appointment.AppointmentWithCustomer, error)
	ScheduleAppointment(bookAppointment appointment.BookAppointmentInput) (appointment.Appointment, error)
	RescheduleAppointment(oldAppointmentId uint, newAppointment uint) (appointment.Appointment, error)
	CancelAppointmentByCustIdAndAppointmentId(custId uint, appointmentId uint) error
	AppointmentResolution(appointmentId uint, employeeId uint, purpose string) error
	ScheduleJoinAccountAppointment(customerId uint, bookAppointment appointment.BookJointAccountAppointmentInput) (appointment.Appointment, error)
}

type appointmentService struct {
}

func (a *appointmentService) ScheduleJoinAccountAppointment(customerInitiated uint, bookAppointment appointment.BookJointAccountAppointmentInput) (appointment.Appointment, error) {

	dbAppointment, err := models.GetAppointmentById(bookAppointment.AppointmentID)

	if err != nil {
		return appointment.Appointment{}, err
	}

	if dbAppointment.Status != "Available" {
		return appointment.Appointment{}, errors.New("appointment is not available")
	}

	customerIds := []uint{}

	for _, customer := range bookAppointment.Customers {
		dbCustomer, err := models.GetCustomerByEmailPhoneAndName(customer.Email, customer.Phone, customer.FirstName, customer.LastName)

		if err != nil {
			return appointment.Appointment{}, errors.New("customer not found " + customer.Email + " " + customer.Phone + " " + customer.FirstName + " " + customer.LastName)
		}

		customerIds = append(customerIds, dbCustomer.ID)
	}

	customerIds = append(customerIds, customerInitiated)

	customerIds = utils.RemoveDuplicatesUint(customerIds)

	customerIdsJson, err := json.Marshal(customerIds)
	if err != nil {
		return appointment.Appointment{}, errors.New("failed to convert customer IDs to JSON")
	}

	dbAppointment.CustomerID = &customerInitiated
	dbAppointment.Status = "Scheduled"
	dbAppointment.Description = "Joint Account Creation"

	dbAppointment, err = models.ScheduleAppointment(dbAppointment)

	if err != nil {
		return appointment.Appointment{}, err
	}

	bookedAppointment, err := models.GetLatestBookedAppointmentByAppointmentIdPurposeAndCustomerId(bookAppointment.AppointmentID, "Joint Account Creation", &customerInitiated)

	jointAccountRequest := models.BankAccountRequest{
		BookedAppointmentId: bookedAppointment.ID,
		CustomerIDs:         string(customerIdsJson),
		AccountType:         "Joint",
		Status:              "Pending",
		Reason:              "Joint Account Creation",
	}

	jointAccountRequest, err = models.CreateBankAccountRequest(jointAccountRequest)

	if err != nil {
		return appointment.Appointment{}, err
	}

	appointmentDTO := convertToAppointmentDTO(dbAppointment)

	return appointmentDTO, nil
}

func (a *appointmentService) AppointmentResolution(appointmentId uint, employeeId uint, purpose string) error {
	dbAppointment, err := models.GetAppointmentById(appointmentId)

	if err != nil {
		return err
	}

	if dbAppointment.Status != "Scheduled" {
		return errors.New("appointment is not scheduled")
	}

	if dbAppointment.EmployeeID != employeeId {
		return errors.New("appointment is not scheduled for this employee")
	}

	if dbAppointment.Description == "Query" && purpose == "Query" {
		_, err := models.ResolveQueryAppointment(appointmentId, employeeId)

		if err != nil {
			return err
		}
	} else if dbAppointment.Description == "Joint Account Creation" && purpose == "jointAccountCreation" {
		_, err := models.ResolveQueryAppointment(appointmentId, employeeId)

		if err != nil {
			return err
		}

		bookedAppointment, err := models.GetLatestBookedAppointmentByAppointmentIdPurposeAndCustomerId(appointmentId, "Joint Account Creation", dbAppointment.CustomerID)

		if err != nil {
			return err
		}

		jointAccountRequest, err := models.GetBankAccountRequestByBookedAppointmentId(bookedAppointment.ID)

		customerIds := []uint{}

		err = json.Unmarshal([]byte(jointAccountRequest.CustomerIDs), &customerIds)

		if err != nil {
			return err
		}

		for _, customerId := range customerIds {
			_, err = models.CreateBankAccount(models.BankAccount{
				CustomerID:     customerId,
				AccountType:    jointAccountRequest.AccountType,
				AccountNumber:  utils.GenerateRandomNumberString(10),
				AccountBalance: 0,
			})

			if err != nil {
				return err
			}

			jointAccountRequest.Status = "Completed"

			_, err := models.UpdateBankAccountRequestStatus(jointAccountRequest.ID, jointAccountRequest.Status, jointAccountRequest.Reason)

			if err != nil {
				return err
			}

			_, err = models.MarkAppointmentAsCompletedByCustIdAndAppointmentId(*dbAppointment.CustomerID, appointmentId)

			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	} else {
		return errors.New("appointment is not scheduled for this purpose")
	}

	return nil
}

func (a *appointmentService) CancelAppointmentByCustIdAndAppointmentId(custId uint, id uint) error {
	_, err := models.CancelAppointmentByCustIdAndAppointmentId(custId, id)

	if err != nil {
		return err
	}

	return nil
}

func (a *appointmentService) RescheduleAppointment(oldAppointmentId uint, newAppointmentId uint) (appointment.Appointment, error) {

	oldAppointment, err := models.GetAppointmentById(oldAppointmentId)

	if err != nil {
		return appointment.Appointment{}, err
	}

	if oldAppointment.Status != "Scheduled" {
		return appointment.Appointment{}, errors.New("appointment is not scheduled")
	}

	newAppointment, err := models.GetAppointmentById(newAppointmentId)

	if err != nil {
		return appointment.Appointment{}, err
	}

	if newAppointment.Status != "Available" {
		return appointment.Appointment{}, errors.New("new appointment is not available")
	}

	rescheduledAppointment, err := models.RescheduleAppointment(oldAppointment, newAppointment)

	if err != nil {
		return appointment.Appointment{}, err
	}

	appointmentDTO := convertToAppointmentDTO(rescheduledAppointment)

	return appointmentDTO, nil
}

func (a *appointmentService) ScheduleAppointment(bookAppointment appointment.BookAppointmentInput) (appointment.Appointment, error) {

	dbAppointment, err := models.GetAppointmentById(bookAppointment.AppointmentID)

	if err != nil {
		return appointment.Appointment{}, err
	}

	if dbAppointment.Status != "Available" {
		return appointment.Appointment{}, errors.New("appointment is not available")
	}

	dbAppointment.CustomerID = bookAppointment.CustomerId
	dbAppointment.Status = "Scheduled"
	dbAppointment.Description = bookAppointment.Purpose

	dbAppointment, err = models.ScheduleAppointment(dbAppointment)

	if err != nil {
		return appointment.Appointment{}, err
	}

	appointmentDTO := convertToAppointmentDTO(dbAppointment)

	return appointmentDTO, nil

}

func (a *appointmentService) GetAppointmentsForEmployee(employeeId uint) ([]appointment.AppointmentWithCustomer, error) {

	dbAppointments, err := models.GetAppointmentsForEmployee(employeeId)

	if err != nil {
		return []appointment.AppointmentWithCustomer{}, err
	}

	var appointments []appointment.AppointmentWithCustomer

	for _, dbAppointment := range dbAppointments {
		appointmentDTO := convertToAppointmentWithCustomerDTO(dbAppointment)
		appointments = append(appointments, appointmentDTO)
	}

	return appointments, nil
}

func (a *appointmentService) GetBookedAppointmentsForCustomer(custId interface{}) ([]appointment.Appointment, error) {

	dbAppointments, err := models.GetBookedAppointmentsForCustomer(custId)

	if err != nil {
		return []appointment.Appointment{}, err
	}

	var appointments []appointment.Appointment

	for _, dbAppointment := range dbAppointments {
		appointmentDTO := convertToAppointmentDTO(dbAppointment)
		appointments = append(appointments, appointmentDTO)
	}

	return appointments, nil
}

func (a *appointmentService) GetAppointmentById(id uint) (appointment.Appointment, error) {

	dbAppointment, err := models.GetAppointmentById(id)

	if err != nil {
		return appointment.Appointment{}, err
	}

	appointmentDTO := convertToAppointmentDTO(dbAppointment)

	return appointmentDTO, nil
}

func (a *appointmentService) GetAppointmentByCustIdAndId(custId interface{}, id uint) (appointment.Appointment, error) {

	dbAppointment, err := models.GetAppointmentByCustIdAndId(custId, id)

	if err != nil {
		return appointment.Appointment{}, err
	}

	appointmentDTO := convertToAppointmentDTO(dbAppointment)

	return appointmentDTO, nil
}

func (a *appointmentService) GetAvailableAppointments(date string) ([]appointment.AvailableAppointment, error) {

	res, err := models.GetAvailableAppointmentsForDate(date)

	if err != nil {
		return []appointment.AvailableAppointment{}, err
	}

	availableAppointments := convertToAvailableAppointmentDTO(res)

	return availableAppointments, nil

}

func (a *appointmentService) CreateAppointment(appointment appointment.Appointment) (appointment.Appointment, error) {

	dbAppointment := convertToAppointmentDB(appointment)

	dbAppointment, err := models.CreateAppointment(dbAppointment)

	if err != nil {
		return appointment, err
	}

	appointmentDTO := convertToAppointmentDTO(dbAppointment)

	return appointmentDTO, nil
}

func convertToAvailableAppointmentDTO(dbAppointments []models.Appointment) []appointment.AvailableAppointment {
	var availableAppointments []appointment.AvailableAppointment

	for _, dbAvailableAppointment := range dbAppointments {
		availableAppointment := appointment.AvailableAppointment{
			StartTime: dbAvailableAppointment.StartTime,
			EndTime:   dbAvailableAppointment.EndTime,
			Date:      dbAvailableAppointment.Date,
			Branch:    dbAvailableAppointment.Branch,
			Status:    dbAvailableAppointment.Status,
			ID:        dbAvailableAppointment.ID,
		}

		availableAppointments = append(availableAppointments, availableAppointment)
	}

	return availableAppointments
}

func convertToAppointmentDB(appointment appointment.Appointment) models.Appointment {
	dbAppointment := models.Appointment{
		EmployeeID:  appointment.EmployeeID,
		CustomerID:  appointment.CustomerID,
		StartTime:   appointment.StartTime,
		EndTime:     appointment.EndTime,
		Date:        appointment.Date,
		Branch:      appointment.Branch,
		Status:      appointment.Status,
		Description: appointment.Description,
	}

	return dbAppointment
}

func convertToAppointmentDTO(dbAppointment models.Appointment) appointment.Appointment {
	appointmentDTO := appointment.Appointment{
		EmployeeID:  dbAppointment.EmployeeID,
		CustomerID:  dbAppointment.CustomerID,
		StartTime:   dbAppointment.StartTime,
		EndTime:     dbAppointment.EndTime,
		Date:        dbAppointment.Date,
		Branch:      dbAppointment.Branch,
		Status:      dbAppointment.Status,
		Description: dbAppointment.Description,
		ID:          dbAppointment.ID,
	}

	return appointmentDTO
}

func convertToAppointmentWithCustomerDTO(dbAppointment models.Appointment) appointment.AppointmentWithCustomer {
	appointmentDTO := appointment.AppointmentWithCustomer{
		EmployeeID:  dbAppointment.EmployeeID,
		StartTime:   dbAppointment.StartTime,
		EndTime:     dbAppointment.EndTime,
		Date:        dbAppointment.Date,
		Branch:      dbAppointment.Branch,
		Status:      dbAppointment.Status,
		Description: dbAppointment.Description,
		ID:          dbAppointment.ID,
		Customer:    customerEntity.Customer{},
	}

	if dbAppointment.CustomerID != nil {
		dbCustomer, err := models.GetCustomerById(*dbAppointment.CustomerID)

		if err != nil {
			return appointmentDTO
		}
	
		customerDTO := customerEntity.Customer{
			CustomerId:  dbCustomer.ID,
			FirstName:   dbCustomer.FirstName,
			LastName:    dbCustomer.LastName,
			Address:     dbCustomer.Address,
			DateOfBirth: dbCustomer.DateOfBirth,
			Email:       dbCustomer.Email,
			Phone:       dbCustomer.Phone,
			UserName:    dbCustomer.UserName,
		}
	
		appointmentDTO.Customer = customerDTO
	}

	return appointmentDTO
}

func NewAppointmentService() AppointmentService {
	return &appointmentService{}
}
