package validator

import (
	"fmt"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
)

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 32); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("must contain only letters, numbers and underscores")
	}

	return nil
}
