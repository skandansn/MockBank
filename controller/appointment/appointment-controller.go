package appointment

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/skandansn/webDevBankBackend/controller/auth"
	"github.com/skandansn/webDevBankBackend/entity/appointment"
	"github.com/skandansn/webDevBankBackend/models"
	appointmentService "github.com/skandansn/webDevBankBackend/service/appointment"
	"github.com/skandansn/webDevBankBackend/utils"
	"log"
	"time"
)

type AppointmentController interface {
	GetAvailableAppointments(ctx *gin.Context) ([]appointment.AvailableAppointment, error)
	CreateAppointment(ctx *gin.Context) (appointment.Appointment, error)
	GetAppointmentById(ctx *gin.Context) (appointment.Appointment, error)
	GetBookedAppointmentsForCustomer(ctx *gin.Context) ([]appointment.Appointment, error)
	GetAppointmentsForEmployee(ctx *gin.Context) ([]appointment.AppointmentWithCustomer, error)
	ScheduleAppointment(ctx *gin.Context) (appointment.Appointment, error)
	RescheduleAppointment(ctx *gin.Context) (appointment.Appointment, error)
	CancelAppointment(ctx *gin.Context) error
	AppointmentResolution(ctx *gin.Context) error
	ScheduleJoinAccountAppointment(ctx *gin.Context) (appointment.Appointment, error)
}

type appointmentController struct {
	service appointmentService.AppointmentService
}

func (c *appointmentController) ScheduleJoinAccountAppointment(ctx *gin.Context) (appointment.Appointment, error) {
	var bjaai appointment.BookJointAccountAppointmentInput

	if err := ctx.ShouldBindJSON(&bjaai); err != nil {
		return appointment.Appointment{}, err
	}

	if len(bjaai.Customers) < 1 {
		return appointment.Appointment{}, errors.New("at least two customers are required for a joint account")
	}

	custId := ctx.GetString("userTypeId")

	appointmentScheduled, err := c.service.ScheduleJoinAccountAppointment(utils.ParseStringAsInt(custId), bjaai)

	if err != nil {
		return appointment.Appointment{}, err
	}

	return appointmentScheduled, nil
}

func (c *appointmentController) AppointmentResolution(ctx *gin.Context) error {
	id := ctx.Param("id")

	if id == "" {
		return errors.New("id is required")
	}

	purpose := ctx.Param("purpose")

	if purpose == "" {
		return errors.New("purpose is required")
	}

	userType := ctx.GetString("userType")
	userTypeId := ctx.GetString("userTypeId")

	if userType != "employee" && userType != "admin" {
		return errors.New("only employees can resolve appointments")
	}

	empId := utils.ParseStringAsInt(userTypeId)

	accessFlag, err := models.IsAccessPresentForEmployee(empId, "create_customer")

	if err != nil {
		return err
	}

	if !accessFlag {
		return errors.New("employee does not have access to resolve appointments")
	}

	idInt := utils.ParseStringAsInt(id)

	err = c.service.AppointmentResolution(idInt, empId, purpose)

	if err != nil {
		return err
	}

	return nil
}

func (c *appointmentController) CancelAppointment(ctx *gin.Context) error {
	id := ctx.Query("id")

	if id == "" {
		return errors.New("id is required")
	}

	userType := ctx.GetString("userType")
	userTypeId := ctx.GetString("userTypeId")

	if userType != "customer" {
		return errors.New("only customers can cancel appointments")
	}

	custId := utils.ParseStringAsInt(userTypeId)

	idInt := utils.ParseStringAsInt(id)

	err := c.service.CancelAppointmentByCustIdAndAppointmentId(custId, idInt)

	if err != nil {
		return err
	}

	return nil
}

func (c *appointmentController) RescheduleAppointment(ctx *gin.Context) (appointment.Appointment, error) {
	var ra appointment.RescheduleAppointmentInput

	if err := ctx.ShouldBindJSON(&ra); err != nil {
		return appointment.Appointment{}, err
	}

	a, err := c.service.RescheduleAppointment(ra.OldAppointmentId, ra.NewAppointmentId)

	if err != nil {
		return appointment.Appointment{}, err
	}

	return a, nil
}

func (c *appointmentController) ScheduleAppointment(ctx *gin.Context) (appointment.Appointment, error) {
	var ba appointment.BookAppointmentInput

	if err := ctx.ShouldBindJSON(&ba); err != nil {
		return appointment.Appointment{}, err
	}

	if ba.CustomerId == nil {
		customer, err := auth.RegisterFromAppointment(ba)

		if err != nil {
			return appointment.Appointment{}, err
		} else {
			ba.CustomerId = &customer.CustomerId
		}
	}

	a, err := c.service.ScheduleAppointment(ba)

	if err != nil {
		return a, err
	}

	return a, nil
}

func (c *appointmentController) GetAppointmentsForEmployee(ctx *gin.Context) ([]appointment.AppointmentWithCustomer, error) {
	userTypeId := ctx.GetString("userTypeId")

	employeeId := utils.ParseStringAsInt(userTypeId)

	appointments, err := c.service.GetAppointmentsForEmployee(employeeId)

	if err != nil {
		return []appointment.AppointmentWithCustomer{}, err
	}

	return appointments, nil
}

func (c *appointmentController) GetBookedAppointmentsForCustomer(ctx *gin.Context) ([]appointment.Appointment, error) {
	userTypeId := ctx.GetString("userTypeId")

	custId := utils.ParseStringAsInt(userTypeId)

	appointments, err := c.service.GetBookedAppointmentsForCustomer(custId)

	if err != nil {
		return []appointment.Appointment{}, err
	}

	return appointments, nil
}

func (c *appointmentController) GetAppointmentById(ctx *gin.Context) (appointment.Appointment, error) {
	id := ctx.Param("id")

	if id == "" {
		return appointment.Appointment{}, nil
	}

	idInt := utils.ParseStringAsInt(id)

	userType := ctx.GetString("userType")
	userTypeId := ctx.GetString("userTypeId")

	if userType == "customer" {
		custId := utils.ParseStringAsInt(userTypeId)
		a, err := c.service.GetAppointmentByCustIdAndId(custId, idInt)

		if err != nil {
			return a, err
		}

		return a, nil
	}

	if userType == "employee" || userType == "admin" {
		a, err := c.service.GetAppointmentById(idInt)

		if err != nil {
			return a, err
		}

		return a, nil
	}

	return appointment.Appointment{}, nil
}

func (c *appointmentController) GetAvailableAppointments(ctx *gin.Context) ([]appointment.AvailableAppointment, error) {
	date := ctx.Query("selectedDate")

	if date == "" {
		date = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	}

	res, err := c.service.GetAvailableAppointments(date)

	if err != nil {
		return []appointment.AvailableAppointment{}, err
	}

	return res, nil
}

func (c *appointmentController) CreateAppointment(ctx *gin.Context) (appointment.Appointment, error) {
	var a appointment.Appointment

	if err := ctx.ShouldBindJSON(&a); err != nil {
		return a, err
	}

	log.Println(a)

	a, err := c.service.CreateAppointment(a)

	if err != nil {
		return a, err
	}

	return a, nil
}

func NewAppointmentController(service appointmentService.AppointmentService) AppointmentController {
	return &appointmentController{
		service: service,
	}
}
