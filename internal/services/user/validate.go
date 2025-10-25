package user

import "strings"

func ValidateEmail(email string) bool {
	email = NormalizeEmail(email)
	if len(email) < 5 || strings.Count(email, "@") != 1 || !strings.Contains(email, ".") {
		return false
	}
	return true
}

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}
