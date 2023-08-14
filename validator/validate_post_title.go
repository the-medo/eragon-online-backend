package validator

func ValidatePostTitle(value string) error {
	if err := ValidateString(value, 3, 256); err != nil {
		return err
	}

	return nil
}
