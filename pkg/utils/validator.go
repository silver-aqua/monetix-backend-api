package utils

import (
	"net/mail"
	"regexp"
	"unicode"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidPhone(phone string) bool {
	// Simple phone validation - adjust based on requirements
	matched, _ := regexp.MatchString(`^\+?[1-9]\d{1,14}$`, phone)
	return matched
}

func ContainsUpper(s string) bool {
	for _, char := range s {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

func ContainsLower(s string) bool {
	for _, char := range s {
		if unicode.IsLower(char) {
			return true
		}
	}
	return false
}

func ContainsDigit(s string) bool {
	for _, char := range s {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}
