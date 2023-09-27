package main

import "net/mail"

// This function is used to validate email.
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// This function is used to validate name.
func ValidateName(name string) bool {
	return len(name) > 3
	return false
}
