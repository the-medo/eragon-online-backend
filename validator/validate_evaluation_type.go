package validator

import "fmt"

func ValidateEvaluationType(value string) error {
	fields := []string{"self", "dm"}

	if StringInSlice(value, fields) == false {
		return fmt.Errorf("incorrect evaluation type! %s", value)
	}

	return nil
}
