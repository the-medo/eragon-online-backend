package validator

import "fmt"

func ValidateInt(value int32, minValue int32, maxValueArgs ...int32) error {
	maxValue := int32(0)

	if len(maxValueArgs) > 0 {
		maxValue = maxValueArgs[0]
	}

	if value < minValue {
		return fmt.Errorf("must be higher or equal than %d", minValue)
	}

	if value > maxValue && len(maxValueArgs) > 0 {
		return fmt.Errorf("must be lower or equal than %d", maxValue)
	}

	return nil
}
