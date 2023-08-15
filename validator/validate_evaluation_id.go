package validator

func ValidateEvaluationId(value int32) error {
	return ValidateInt(value, 1)
}
