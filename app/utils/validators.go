package utils

import (
	"regexp"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

// La contraseña debe contener al menos 8 caracteres
var PasswordValidator validator.Func = func(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return utf8.RuneCountInString(password) >= 8
}

// El número de teléfono debe contener solo dígitos
var PhoneValidator validator.Func = func(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return isAllDigits(phone)
}

// La placa debe cumplir con el formato estándar en Guatemala
var LicensePlateValidator validator.Func = func(fl validator.FieldLevel) bool {
	plate := fl.Field().String()
	re := regexp.MustCompile(`^[A-Za-z]\d{3}[A-Z]{3}$`)
	return re.MatchString(plate)
}

// El texto de una pregunta debe ser no vacío y contener menos de 51 caracteres
var QuestionTextValidator validator.Func = func(fl validator.FieldLevel) bool {
	text := fl.Field().String()
	length := utf8.RuneCountInString(text)
	return length > 0 && length <= 50
}

// La descripción debe contener máximo 255 caracteres
var QuestionDescriptionValidator validator.Func = func(fl validator.FieldLevel) bool {
	desc := fl.Field().String()
	return utf8.RuneCountInString(desc) <= 255
}

// El tipo debe contener entre 1 y 50 caracteres
var QuestionTypeValidator validator.Func = func(fl validator.FieldLevel) bool {
	typ := fl.Field().String()
	length := utf8.RuneCountInString(typ)
	return length > 0 && length <= 50
}

// La opción debe contener entre 1 y 50 caracteres
var QuestionOptionValidator validator.Func = func(fl validator.FieldLevel) bool {
	opt := fl.Field().String()
	length := utf8.RuneCountInString(opt)
	return length > 0 && length <= 50
}

// El estado debe pertenecer a una lista válida
var SubmissionStatusValidator validator.Func = func(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	return status == "pending" || status == "approved" || status == "rejected"
}
