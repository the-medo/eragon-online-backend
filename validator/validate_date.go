package validator

import (
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type DateBound struct {
	Min time.Time
	Max time.Time
}

// ValidateDate checks if a protobuf timestamp is between a minimum and a maximum time.
func ValidateDate(date *timestamppb.Timestamp, bound DateBound) error {
	t := date.AsTime()

	if !bound.Min.IsZero() && t.Before(bound.Min) {
		return fmt.Errorf("date must be after %s", bound.Min)
	}

	if !bound.Max.IsZero() && t.After(bound.Max) {
		return fmt.Errorf("date must be before %s", bound.Max)
	}

	if bound.Min.IsZero() && bound.Max.IsZero() {
		return errors.New("at least one time bound (min or max) must be provided")
	}

	return nil
}
