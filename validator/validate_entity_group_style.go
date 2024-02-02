package validator

import (
	"fmt"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/pb"
)

type EntityGroupStyle string

// ValidateEntityGroupStyle checks if the provided string matches any of the EntityGroupStyle values
func ValidateEntityGroupStyle(value string) error {
	str := converters.ConvertEntityGroupStyleToPB(value)

	if str == pb.EntityGroupStyle_ENTITY_GROUP_STYLE_UNKNOWN {
		return fmt.Errorf("invalid entity group style: %s", value)
	}
	return nil
}
