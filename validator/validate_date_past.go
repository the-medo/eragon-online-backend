package validator

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// ValidateDatePast ensures the provided protobuf timestamp is in the past.
func ValidateDatePast(date *timestamppb.Timestamp) error {
	now := time.Now()
	bound := DateBound{
		Max: now,
	}

	return ValidateDate(date, bound)
}
