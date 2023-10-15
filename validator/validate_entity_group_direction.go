package validator

import "fmt"

type EntityGroupDirection string

const (
	Horizontal EntityGroupDirection = "horizontal"
	Vertical   EntityGroupDirection = "vertical"
)

// ValidateEntityGroupDirection checks if the provided string matches any of the EntityGroupDirection values
func ValidateEntityGroupDirection(value string) error {
	egs := EntityGroupDirection(value)
	switch egs {
	case Horizontal, Vertical:
		return nil
	default:
		return fmt.Errorf("invalid entity group direction: %s", value)
	}
}
