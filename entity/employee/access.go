package employee

type AccessUpdateInput struct {
	AccessList []Access `json:"accessList" binding:"required"`
}

type Access struct {
	AccessName  string `json:"accessName" binding:"required"`
	AccessGiven bool   `json:"accessGiven" binding:"required"`
}

func GetEmployeeAccessItems() map[string]bool {
	accessItems := make(map[string]bool)
	accessItems["view_customer_details"] = true
	accessItems["approve_card"] = true
	accessItems["create_customer"] = true
	accessItems["view_customer_transactions"] = true
	return accessItems
}
