package utils

import (
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

// RoleValidator checks if the role string is valid.
// Valid roles are "admin" or "driver".
var RoleValidator validator.Func = func(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	return role == "admin" || role == "driver"
}

// PasswordValidator validates the password field.
// The password must be at least 8 characters long.
var PasswordValidator validator.Func = func(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return utf8.RuneCountInString(password) >= 8
}

// PhoneValidator validates the phone number field.
// The phone number must contain only digits.
var PhoneValidator validator.Func = func(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return IsAllDigits(phone)
}