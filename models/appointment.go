package models

import "github.com/jinzhu/gorm"

type Appointment struct {
	gorm.Model
	EmployeeID  uint   `gorm:"not null" json:"employee_id"`
	CustomerID  *uint  `gorm:"" json:"customer_id"`
	StartTime   string `gorm:"size:255;not null;" json:"startTime"`
	EndTime     string `gorm:"size:255;not null;" json:"endTime"`
	Date        string `gorm:"size:255;not null;" json:"date"`
	Branch      string `gorm:"size:255;not null;" json:"branch"`
	Status      string `gorm:"size:255;not null;" json:"status"`
	Description string `gorm:"size:255;not null;" json:"description"`
}

type BookedAppointment struct {
	gorm.Model
	EmployeeID    uint        `gorm:"not null" json:"employee_id"`
	CustomerID    uint        `gorm:"not null" json:"customer_id"`
	AppointmentID uint        `gorm:"not null" json:"appointment_id"`
	Appointment   Appointment `gorm:"" json:"appointment"`
	Status        string      `gorm:"size:255;not null;" json:"status"`
	Purpose       string      `gorm:"size:255;not null;" json:"purpose"`
}

func ResolveQueryAppointment(id uint, employeeId uint) (Appointment, error) {

	var a Appointment

	if err := DB.Where("id = ? AND employee_id = ? AND status = ?", id, employeeId, "Scheduled").First(&a).Error; err != nil {
		return a, err
	}

	a.Status = "Completed"

	err := DB.Save(&a).Error

	if err != nil {
		return a, err
	}

	var ba BookedAppointment

	if err := DB.Where("appointment_id = ?", id).First(&ba).Error; err != nil {
		return a, err
	}

	ba.Status = "Completed"

	err = DB.Save(&ba).Error

	return a, nil
}

func GetAvailableAppointmentsForDate(date string) ([]Appointment, error) {

	var a []Appointment

	if err := DB.Where("date = ? AND status = ?", date, "Available").Find(&a).Error; err != nil {
		return a, err
	}

	return a, nil
}

func CreateAppointment(a Appointment) (Appointment, error) {

	err := DB.Create(&a).Error
	if err != nil {
		return Appointment{}, err
	}

	return a, nil
}

func GetLatestBookedAppointmentByAppointmentIdPurposeAndCustomerId(appointmentId uint, purpose string, customerId *uint) (BookedAppointment, error) {

	var ba BookedAppointment

	if err := DB.Where("appointment_id = ? AND purpose = ? AND customer_id = ?", appointmentId, purpose, customerId).Last(&ba).Error; err != nil {
		return ba, err
	}

	return ba, nil
}

func GetAppointmentById(id uint) (Appointment, error) {

	var a Appointment

	if err := DB.Where("id = ?", id).First(&a).Error; err != nil {
		return a, err
	}

	return a, nil
}

func GetAppointmentByCustIdAndId(custId interface{}, id uint) (Appointment, error) {

	var a Appointment

	if err := DB.Where("id = ? AND customer_id = ?", id, custId).First(&a).Error; err != nil {
		return a, err
	}

	return a, nil
}

func GetBookedAppointmentsForCustomer(custId interface{}) ([]Appointment, error) {

	var a []Appointment

	if err := DB.Where("customer_id = ? AND status != ?", custId, "Available").Find(&a).Error; err != nil {
		return a, err
	}

	return a, nil
}

func GetAppointmentsForEmployee(employeeId uint) ([]Appointment, error) {

	var a []Appointment

	if err := DB.Where("employee_id = ?", employeeId).Find(&a).Error; err != nil {
		return a, err
	}

	return a, nil
}

func ScheduleAppointment(a Appointment) (Appointment, error) {

	ba := BookedAppointment{
		EmployeeID:    a.EmployeeID,
		CustomerID:    *a.CustomerID,
		AppointmentID: a.ID,
		Status:        "Scheduled",
		Purpose:       a.Description,
	}

	err := DB.Create(&ba).Error
	if err != nil {
		return Appointment{}, err
	}

	a.Status = "Scheduled"
	a.CustomerID = &ba.CustomerID
	a.Description = ba.Purpose
	err = DB.Save(&a).Error

	if err != nil {
		return Appointment{}, err
	}

	return a, nil
}

func RescheduleAppointment(old Appointment, new Appointment) (Appointment, error) {

	old.Status = "Available"
	old.CustomerID = nil
	err := DB.Save(&old).Error

	if err != nil {
		return Appointment{}, err
	}

	oldBooked := BookedAppointment{}
	err = DB.Where("appointment_id = ?", old.ID).First(&oldBooked).Error

	if err != nil {
		return Appointment{}, err
	}

	oldBooked.Status = "Cancelled"

	err = DB.Save(&oldBooked).Error

	if err != nil {
		return Appointment{}, err
	}

	new.Status = "Scheduled"
	new.CustomerID = &oldBooked.CustomerID
	err = DB.Save(&new).Error
	if err != nil {
		return Appointment{}, err
	}

	newBooked := BookedAppointment{
		EmployeeID:    new.EmployeeID,
		CustomerID:    *new.CustomerID,
		AppointmentID: new.ID,
		Status:        "Scheduled",
	}

	err = DB.Create(&newBooked).Error

	if err != nil {
		return Appointment{}, err
	}

	return new, nil
}

func CancelAppointmentByCustIdAndAppointmentId(custId uint, id uint) (Appointment, error) {

	a := Appointment{}

	err := DB.Where("id = ? AND customer_id = ?", id, custId).First(&a).Error

	if err != nil {
		return Appointment{}, err
	}

	a.Status = "Available"
	a.CustomerID = nil

	err = DB.Save(&a).Error

	if err != nil {
		return Appointment{}, err
	}

	ba := BookedAppointment{}

	err = DB.Where("appointment_id = ? AND customer_id = ?", id, custId).First(&ba).Error

	if err != nil {
		return Appointment{}, err
	}

	ba.Status = "Cancelled"

	err = DB.Save(&ba).Error

	if err != nil {
		return Appointment{}, err
	}

	return a, nil
}

func MarkAppointmentAsCompletedByCustIdAndAppointmentId(custId uint, id uint) (Appointment, error) {

	a := Appointment{}

	err := DB.Where("id = ? AND customer_id = ?", id, custId).First(&a).Error

	if err != nil {
		return Appointment{}, err
	}

	a.Status = "Completed"

	err = DB.Save(&a).Error

	if err != nil {
		return Appointment{}, err
	}

	ba := BookedAppointment{}

	err = DB.Where("appointment_id = ? AND customer_id = ?", id, custId).First(&ba).Error

	if err != nil {
		return Appointment{}, err
	}

	ba.Status = "Completed"

	err = DB.Save(&ba).Error

	if err != nil {
		return Appointment{}, err
	}

	return a, nil
}
