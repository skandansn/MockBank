package routes

import (
	"github.com/skandansn/webDevBankBackend/controller"
	accountController "github.com/skandansn/webDevBankBackend/controller/account"
	appointmentController "github.com/skandansn/webDevBankBackend/controller/appointment"
	"github.com/skandansn/webDevBankBackend/controller/auth"
	cardController "github.com/skandansn/webDevBankBackend/controller/card"
	customerController "github.com/skandansn/webDevBankBackend/controller/customer"
	employeeController "github.com/skandansn/webDevBankBackend/controller/employee"
	transactionController "github.com/skandansn/webDevBankBackend/controller/transaction"
	"github.com/skandansn/webDevBankBackend/service"
	accountService "github.com/skandansn/webDevBankBackend/service/account"
	appointmentService "github.com/skandansn/webDevBankBackend/service/appointment"
	cardService "github.com/skandansn/webDevBankBackend/service/card"
	customerService "github.com/skandansn/webDevBankBackend/service/customer"
	employeeService "github.com/skandansn/webDevBankBackend/service/employee"
	transactionService "github.com/skandansn/webDevBankBackend/service/transaction"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	appointmentServiceInstance    = appointmentService.NewAppointmentService()
	appointmentControllerInstance = appointmentController.NewAppointmentController(appointmentServiceInstance)
	employeeServiceInstance       = employeeService.New()
	employeeControllerInstance    = employeeController.New(employeeServiceInstance)

	employeeAccessServiceInstance    = employeeService.NewAccessService()
	employeeAccessControllerInstance = employeeController.NewAccessController(employeeAccessServiceInstance)

	customerServiceInstance    = customerService.NewCustomerService()
	customerControllerInstance = customerController.NewCustomerController(customerServiceInstance)

	bankAccountServiceInstance    = accountService.NewBankAccountService()
	bankAccountControllerInstance = accountController.NewBankAccountController(bankAccountServiceInstance)

	cardServiceInstance    = cardService.NewCardService()
	cardControllerInstance = cardController.NewCardController(cardServiceInstance)

	transactionServiceInstance    = transactionService.NewTransactionService()
	transactionControllerInstance = transactionController.NewTransactionController(transactionServiceInstance)

	companyBranchLocationServiceInstance    = service.NewCompanyBranchLocationService()
	companyBranchLocationControllerInstance = controller.NewCompanyBranchLocationController(companyBranchLocationServiceInstance)
)

var Routes = []Route{
	{
		Path:   "/healthCheck",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"public": true, "employee": true, "admin": true, "customer": true},
		Handler: func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "alive!",
			})
		},
	},
	{
		Path:   "/register",
		Method: http.MethodPost,
		Tiers:  map[string]bool{"public": true},
		Handler: func(ctx *gin.Context) {
			cus, err := auth.Register(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(201, cus)
			}
		},
	},
	{
		Path:    "/login",
		Method:  http.MethodPost,
		Tiers:   map[string]bool{"public": true},
		Handler: auth.Login,
	},
	{
		Path:   "/employees",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"admin": true},
		Handler: func(ctx *gin.Context) {
			res, err := employeeControllerInstance.FindAll()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, res)
			}
		},
	},
	{
		Path:   "/employees",
		Method: http.MethodPost,
		Tiers:  map[string]bool{"admin": true},
		Handler: func(ctx *gin.Context) {
			_, err := employeeControllerInstance.Save(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, gin.H{
					"message": "employee created",
				})
			}
		},
	},
	{
		Path:   "/employees",
		Method: http.MethodPut,
		Tiers:  map[string]bool{"admin": true},
		Handler: func(ctx *gin.Context) {
			_, err := employeeControllerInstance.Edit(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, gin.H{
					"message": "employee updated",
				})
			}
		},
	},
	{
		Path:   "/employees/:employeeId",
		Method: http.MethodDelete,
		Tiers:  map[string]bool{"admin": true},
		Handler: func(ctx *gin.Context) {
			err := employeeControllerInstance.Delete(ctx.Param("employeeId"))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, gin.H{
					"message": "employee deleted",
				})
			}
		},
	},
	// Access routes
	{
		Path:   "/employees/:employeeId/access",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"admin": true, "employee": true},
		Handler: func(ctx *gin.Context) {
			res, err := employeeAccessControllerInstance.GetAccessForEmployee(ctx.Param("employeeId"))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, res)
			}
		},
	},
	{
		Path:   "/employees/:employeeId/access",
		Method: http.MethodPut,
		Tiers:  map[string]bool{"admin": true},
		Handler: func(ctx *gin.Context) {
			_, err := employeeAccessControllerInstance.SaveAccessForEmployee(ctx.Param("employeeId"), ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, gin.H{
					"message": "access updated",
				})
			}
		},
	},
	// Customer routes
	{
		Path:   "/customer",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"admin": true, "employee": true, "customer": true},
		Handler: func(ctx *gin.Context) {
			res, err := customerControllerInstance.GetCustomerDetails(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, res)
			}
		},
	},
	{
		Path:   "/customers",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"admin": true, "employee": true},
		Handler: func(ctx *gin.Context) {
			res, err := customerControllerInstance.GetAllCustomers(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, res)
			}
		},
	},
	{
		Path:   "/updateProfile",
		Method: http.MethodPut,
		Tiers:  map[string]bool{"admin": true, "employee": true, "customer": true},
		Handler: func(ctx *gin.Context) {
			res, err := customerControllerInstance.UpdateProfile(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, res)
			}
		},
	},
	{
		Path:    "/user",
		Method:  http.MethodGet,
		Tiers:   map[string]bool{"admin": true, "employee": true, "customer": true},
		Handler: auth.CurrentUser,
	},
	// appointments
	{
		Path:   "/availableAppointments",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"admin": true, "employee": true, "customer": true, "public": true},
		Handler: func(ctx *gin.Context) {
			res, err := appointmentControllerInstance.GetAvailableAppointments(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, res)
			}
		},
	},
	{
		Path:   "/appointments",
		Method: http.MethodPost,
		Tiers:  map[string]bool{"admin": true, "employee": true},
		Handler: func(ctx *gin.Context) {
			res, err := appointmentControllerInstance.CreateAppointment(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, res)
			}
		},
	},
	{
		Path:   "/bookAppointment",
		Method: http.MethodPost,
		Tiers:  map[string]bool{"customer": true, "public": true},
		Handler: func(ctx *gin.Context) {
			res, err := appointmentControllerInstance.ScheduleAppointment(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(201, res)
			}
		},
	},
	{
		Path:   "/bookJoinAccountAppointment",
		Method: http.MethodPost,
		Tiers:  map[string]bool{"customer": true},
		Handler: func(ctx *gin.Context) {
			res, err := appointmentControllerInstance.ScheduleJoinAccountAppointment(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(201, res)
			}
		},
	},
	{
		Path:   "/rescheduleAppointment",
		Method: http.MethodPost,
		Tiers:  map[string]bool{"customer": true},
		Handler: func(ctx *gin.Context) {
			res, err := appointmentControllerInstance.RescheduleAppointment(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(201, res)
			}
		},
	},
	{
		Path:   "/cancelAppointment",
		Method: http.MethodPost,
		Tiers:  map[string]bool{"customer": true},
		Handler: func(ctx *gin.Context) {
			err := appointmentControllerInstance.CancelAppointment(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, "Appointment cancelled")
			}
		},
	},
	{
		Path:   "/employee/appointmentResolution/:purpose/:id",
		Method: http.MethodPost,
		Tiers:  map[string]bool{"employee": true},
		Handler: func(ctx *gin.Context) {
			err := appointmentControllerInstance.AppointmentResolution(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, "Appointment resolved")
			}
		},
	},
	{
		Path:   "/employee/viewAppointments",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"employee": true},
		Handler: func(ctx *gin.Context) {
			res, err := appointmentControllerInstance.GetAppointmentsForEmployee(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, res)
			}
		},
	},
	{
		Path:   "/viewAppointments",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"customer": true},
		Handler: func(ctx *gin.Context) {
			res, err := appointmentControllerInstance.GetBookedAppointmentsForCustomer(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, res)
			}
		},
	},
	{
		Path:   "/viewAppointment/:id",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"admin": true, "employee": true, "customer": true},
		Handler: func(ctx *gin.Context) {
			res, err := appointmentControllerInstance.GetAppointmentById(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, res)
			}
		},
	},
	// Customer focussed routes,
	{
		Path:   "/createCustomer",
		Method: http.MethodPost,
		Tiers:  map[string]bool{"admin": true, "employee": true},
		Handler: func(ctx *gin.Context) {
			res, err := customerControllerInstance.CreateCustomer(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(201, res)
			}
		},
	},
	// Bank account routes
	{
		Path:   "/bankAccounts",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"customer": true},
		Handler: func(ctx *gin.Context) {
			res, err := bankAccountControllerInstance.GetBankAccountsByCustomerId(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, res)
			}
		},
	},
	// card routes
	{
		Path:   "/cards",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"customer": true},
		Handler: func(ctx *gin.Context) {
			res, err := cardControllerInstance.GetCardsByCustomerId(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, res)
			}
		},
	},
	{
		Path:   "/applyCreditCard",
		Method: http.MethodPost,
		Tiers:  map[string]bool{"customer": true},
		Handler: func(ctx *gin.Context) {
			res, err := cardControllerInstance.CreateCardRequest(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusCreated, res)
			}
		},
	},
	{
		Path:   "/viewCardRequests",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"employee": true, "customer": true},
		Handler: func(ctx *gin.Context) {
			res, err := cardControllerInstance.GetPendingCardRequests(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, res)
			}
		},
	},
	{
		Path:   "/approveOrRejectCardRequest",
		Method: http.MethodPost,
		Tiers:  map[string]bool{"employee": true},
		Handler: func(ctx *gin.Context) {
			res, err := cardControllerInstance.ApproveOrRejectCardRequest(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusCreated, res)
			}
		},
	},
	// Transaction routes
	{
		Path:   "/transferMoney",
		Method: http.MethodPost,
		Tiers:  map[string]bool{"customer": true},
		Handler: func(ctx *gin.Context) {
			res, err := transactionControllerInstance.CreateTransaction(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusCreated, res)
			}
		},
	},
	{
		Path:   "/accountTransactions",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"customer": true},
		Handler: func(ctx *gin.Context) {
			res, err := transactionControllerInstance.GetTransactionsForSenderAccount(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, res)
			}
		},
	},
	{
		Path:   "/customerTransactions",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"customer": true},
		Handler: func(ctx *gin.Context) {
			res, err := transactionControllerInstance.GetTransactionsForCustomer(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, res)
			}
		},
	},
	{
		Path:   "/customerTransactionsByEmployee/:customerId",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"employee": true},
		Handler: func(ctx *gin.Context) {
			res, err := transactionControllerInstance.GetTransactionsForCustomerByEmployee(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, res)
			}
		},
	},
	// branches
	{
		Path:   "/branches",
		Method: http.MethodGet,
		Tiers:  map[string]bool{"public": true, "customer": true, "employee": true, "admin": true},
		Handler: func(ctx *gin.Context) {
			res, err := companyBranchLocationControllerInstance.GetCompanyBranchLocations(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, res)
			}
		},
	},
	{
		Path:   "/branches",
		Method: http.MethodPost,
		Tiers:  map[string]bool{"admin": true},
		Handler: func(ctx *gin.Context) {
			res, err := companyBranchLocationControllerInstance.SaveCompanyBranchLocation(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusCreated, res)
			}
		},
	},
}

type Route struct {
	Path    string
	Method  string
	Tiers   map[string]bool
	Handler gin.HandlerFunc
}
