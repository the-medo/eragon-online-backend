package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidName     = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
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

	if !isValidUsername(value) {
		return fmt.Errorf("must contain only lowercase letters, numbers, and underscores")
	}

	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

// ValidateLimitOrOffset value = 0 means no limit
func ValidateLimitOrOffset(value int32, maxValueArgs ...int32) error {
	maxValue := int32(0)

	if len(maxValueArgs) > 0 {
		maxValue = maxValueArgs[0]
	}

	if value < 0 {
		return fmt.Errorf("must be a positive integer or zero")
	}

	if maxValue > 0 {
		if value > maxValue {
			return fmt.Errorf("must be lower or equal to %d", maxValue)
		}
	}
	return nil
}

func ValidateNumber(value int32, minValue int32, maxValue int32) error {
	if value < minValue {
		return fmt.Errorf("must be higher than %d", minValue)
	}

	if value > maxValue {
		return fmt.Errorf("must be lower than %d", maxValue)
	}

	return nil
}

func ValidateUserId(value int32) error {
	if value < 1 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}

func ValidateImgId(value int32) error {
	if value < 1 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}

func ValidateEmailId(value int64) error {
	if value < 1 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}

func ValidateSecretCode(value string) error {
	if err := ValidateString(value, 32, 128); err != nil {
		return err
	}
	return nil
}
