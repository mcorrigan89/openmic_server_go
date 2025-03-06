package validation

import (
	validator "github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	// Register custom validations
	v.RegisterValidation("password_strength", validatePasswordStrength)

	return &Validator{
		validate: v,
	}
}

func validatePasswordStrength(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// Your password strength logic here
	return len(password) >= 8
}
