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
		return "El teléfono debe contener únicamente dígitos"

	case "password":
		return "La contraseña debe contener al menos 8 dígitos"

	case "validrole":
		return "Debe ingresar un rol válido"

	case "email":
		return "Correo inválido"
	}

	return ""
}
