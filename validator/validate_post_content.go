package validator

func ValidatePostContent(value string) error {
	if err := ValidateString(value, 0, 100000); err != nil {
		return err
	}

	return nil
}
