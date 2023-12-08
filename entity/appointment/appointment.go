package appointment

import (
	"github.com/skandansn/webDevBankBackend/entity/auth"
)

type AvailableAppointment struct {
	ID        uint   `gorm:"primary_key;auto_increment" json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Date      string `json:"date"`
	Branch    string `json:"branch"`
	Status    string `json:"status"`
}

type Appointment struct {
	ID          uint   `gorm:"primary_key;auto_increment" json:"id" binding:"required"`
	EmployeeID  uint   `gorm:"not null" json:"employee_id" binding:"required"`
	CustomerID  *uint  `gorm:"not null" json:"customer_id"`
	StartTime   string `gorm:"size:255;not null;" json:"startTime" binding:"required"`
	EndTime     string `gorm:"size:255;not null;" json:"endTime" binding:"required"`
	Date        string `gorm:"size:255;not null;" json:"date" binding:"required"`
	Branch      string `gorm:"size:255;not null;" json:"branch" binding:"required"`
	Status      string `gorm:"size:255;not null;" json:"status" binding:"required"`
	Description string `gorm:"size:255;not null;" json:"description"`
}

type BookAppointmentInput struct {
	AppointmentID   uint                        `json:"appointmentId" binding:"required"`
	CustomerId      *uint                       `json:"customerId" binding:""`
	CustomerDetails *auth.CustomerRegisterInput `json:"customerDetails" binding:""`
}

type RescheduleAppointmentInput struct {
	OldAppointmentId uint `json:"oldAppointmentId" binding:"required"`
	NewAppointmentId uint `json:"newAppointmentId" binding:"required"`
}
