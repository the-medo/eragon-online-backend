package validator

import (
	"fmt"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

type DeleteEntityGroupContentAction string

// ValidateDeleteEntityGroupContentAction checks if the provided string matches any of the DeleteEntityGroupContentAction values
func ValidateDeleteEntityGroupContentAction(value pb.DeleteEntityGroupContentAction) error {
	str := converters.ConvertDeleteEntityGroupContentActionToDB(value)

	if str == db.DeleteEntityGroupContentActionUnknown {
		return fmt.Errorf("invalid DeleteEntityGroupContentAction: %s", value)
	}
	return nil
}
