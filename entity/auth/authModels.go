package auth

type CustomerRegisterInput struct {
	UserName    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastname" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	DateOfBirth string `json:"dateOfBirth" binding:"required"`
	Address     string `json:"address" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
