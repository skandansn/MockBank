package validators

import "github.com/go-playground/validator/v10"

func ValidateEmployeeRole(field validator.FieldLevel) bool {
	return field.Field().String() == "admin" || field.Field().String() == "employee" || field.Field().String() == "manager"
}
