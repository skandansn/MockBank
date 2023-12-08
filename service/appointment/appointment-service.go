package appointment

import (
	"errors"
	"github.com/skandansn/webDevBankBackend/entity/appointment"
	"github.com/skandansn/webDevBankBackend/models"
)

type AppointmentService interface {
	GetAvailableAppointments(date string) ([]appointment.AvailableAppointment, error)
	CreateAppointment(appointment appointment.Appointment) (appointment.Appointment, error)
	GetAppointmentById(id uint) (appointment.Appointment, error)
	GetAppointmentByCustIdAndId(custId interface{}, id uint) (appointment.Appointment, error)
	GetBookedAppointmentsForCustomer(custId interface{}) ([]appointment.Appointment, error)
	GetAppointmentsForEmployee(employeeId uint) ([]appointment.Appointment, error)
	ScheduleAppointment(bookAppointment appointment.BookAppointmentInput) (appointment.Appointment, error)
	RescheduleAppointment(oldAppointmentId uint, newAppointment uint) (appointment.Appointment, error)
	CancelAppointmentByCustIdAndAppointmentId(custId uint, appointmentId uint) error
}

type appointmentService struct {
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

	dbAppointment, err = models.ScheduleAppointment(dbAppointment)

	if err != nil {
		return appointment.Appointment{}, err
	}

	appointmentDTO := convertToAppointmentDTO(dbAppointment)

	return appointmentDTO, nil

}

func (a *appointmentService) GetAppointmentsForEmployee(employeeId uint) ([]appointment.Appointment, error) {

	dbAppointments, err := models.GetAppointmentsForEmployee(employeeId)

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

func NewAppointmentService() AppointmentService {
	return &appointmentService{}
}
