package validator

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// ValidateDateFuture ensures the provided protobuf timestamp is in the past.
func ValidateDateFuture(date *timestamppb.Timestamp) error {
	now := time.Now()
	bound := DateBound{
		Min: now,
	}

	return ValidateDate(date, bound)
}
