package utils

// Determina si una string contiene únicamente dígitos o no
// Utiliza el valor ASCII de cada caracter
func IsAllDigits(s string) bool {
	for _, c := range s {
		if c < 48 || c > 57 {
			return false
		}
	}

	return true
}

func GetTagMessage(tag string) string {
	switch tag {
	case "phone":
		return "Phone must only contain digits"

	case "password":
		return "Password must have at least 8 characters"

	case "validrole":
		return "Must enter a valid role"

	case "email":
		return "Invalid email"
	}

	return ""
}
