package validator

import "fmt"

func ValidateInt64(value int64, minValue int64, maxValueArgs ...int64) error {
	maxValue := int64(0)

	if len(maxValueArgs) > 0 {
		maxValue = maxValueArgs[0]
	}

	if value < minValue {
		return fmt.Errorf("must be higher than %d", minValue)
	}

	if value > maxValue && len(maxValueArgs) > 0 {
		return fmt.Errorf("must be lower or equal than %d", maxValue)
	}

	return nil
}
