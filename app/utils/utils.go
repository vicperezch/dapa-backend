package utils

// IsAllDigits checks if a string contains only digit characters.
// It iterates through each character and verifies its ASCII value.
func IsAllDigits(s string) bool {
	for _, c := range s {
		if c < 48 || c > 57 {
			return false
		}
	}
	return true
}

// GetTagMessage returns a user-friendly error message
// corresponding to a validation tag.
func GetTagMessage(tag string) string {
	switch tag {
	case "phone":
		return "Phone number must contain digits only"

	case "password":
		return "Password must be at least 8 characters long"

	case "email":
		return "Invalid email address"

	case "plate":
		return "Invalid license plate format"
	}

	return "Invalid request format"
}
