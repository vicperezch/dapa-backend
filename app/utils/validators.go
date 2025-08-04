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

// QuestionTextValidator checks that question text is not empty and max 50 characters.
var QuestionTextValidator validator.Func = func(fl validator.FieldLevel) bool {
	text := fl.Field().String()
	length := utf8.RuneCountInString(text)
	return length > 0 && length <= 50
}

// QuestionDescriptionValidator allows optional description up to 255 characters.
var QuestionDescriptionValidator validator.Func = func(fl validator.FieldLevel) bool {
	desc := fl.Field().String()
	return utf8.RuneCountInString(desc) <= 255
}

// QuestionTypeValidator ensures type name is not empty and max 50 characters.
var QuestionTypeValidator validator.Func = func(fl validator.FieldLevel) bool {
	typ := fl.Field().String()
	length := utf8.RuneCountInString(typ)
	return length > 0 && length <= 50
}

// QuestionOptionValidator ensures the option text is not empty and max 50 characters.
var QuestionOptionValidator validator.Func = func(fl validator.FieldLevel) bool {
	opt := fl.Field().String()
	length := utf8.RuneCountInString(opt)
	return length > 0 && length <= 50
}

// SubmissionStatusValidator validates status is one of the allowed values.
var SubmissionStatusValidator validator.Func = func(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	return status == "pending" || status == "approved" || status == "rejected"
}