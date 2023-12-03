package validator

import "fmt"

type EntityGroupStyle string

const (
	Framed    EntityGroupStyle = "framed"
	NonFramed EntityGroupStyle = "non-framed"
)

// ValidateEntityGroupStyle checks if the provided string matches any of the EntityGroupStyle values
func ValidateEntityGroupStyle(value string) error {
	egs := EntityGroupStyle(value)
	switch egs {
	case Framed, NonFramed:
		return nil
	default:
		return fmt.Errorf("invalid entity group style: %s", value)
	}
}
