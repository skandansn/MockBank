package employee

type Employee struct {
	FirstName  string `json:"firstName" binding:"required,min=2,max=50"`
	LastName   string `json:"lastName" binding:"required,min=2,max=50"`
	Email      string `json:"email" binding:"required,email"`
	Phone      string `json:"phone" binding:"required"`
	Role       string `json:"role"  binding:"required,min=2,max=50" validate:"isValidRole"`
	Branch     string `json:"branch"  binding:"required,min=2,max=50"`
	EmployeeID uint   `json:"employeeID"`
	//Manager    *Employee `json:"manager"`
	AccessList []Access `json:"accessList"`
}

type EmployeeRegisterInput struct {
	FirstName  string    `json:"firstName" binding:"required,min=2,max=50"`
	LastName   string    `json:"lastName" binding:"required,min=2,max=50"`
	Email      string    `json:"email" binding:"required,email"`
	Phone      string    `json:"phone" binding:"required"`
	Role       string    `json:"role"  binding:"required,min=2,max=50" validate:"isValidRole"`
	Branch     string    `json:"branch"  binding:"required,min=2,max=50"`
	EmployeeID uint      `json:"employeeID"`
	Manager    *Employee `json:"manager"`
	Username   string    `json:"username" binding:"required,min=2,max=50"`
	Password   string    `json:"password" binding:"required,min=2,max=50"`
	AccessList []Access  `json:"accessList"`
}
