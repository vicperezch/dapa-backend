package utils

import (
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

// Determina si un rol es válido (admin o driver).
var RoleValidator validator.Func = func(fl validator.FieldLevel) bool {
	role := fl.Field().String()

	return role == "admin" || role == "driver"
}

// Realiza validaciones de contraseña.
// Debe contener al menos 8 caracteres.
var PasswordValidator validator.Func = func(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	return utf8.RuneCountInString(password) >= 8
}

// Realiza validaciones para números de teléfono.
// Deben contener únicamente dígitos.
var PhoneValidator validator.Func = func(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	return IsAllDigits(phone)
}
