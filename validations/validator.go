package validations

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullnameRegex = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLength, maxLength int) error {

	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("length must be between %d and %d", minLength, maxLength)
	}

	return nil
}

func ValidateUsername(value string) error {

	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsernameRegex(value) {
		return fmt.Errorf("username must be alphanumeric and contain only underscores")
	}

	return nil
}

func ValidatePassword(value string) error {
	if err := ValidateString(value, 6, 100); err != nil {
		return err
	}
	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("invalid email address")
	}

	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidFullnameRegex(value) {
		return fmt.Errorf("fullname must be letters and spaces")
	}

	return nil
}
