package validator

import "google.golang.org/genproto/googleapis/rpc/errdetails"

func ValidateSliceOfIds(ids []int32, fieldName string) []*errdetails.BadRequest_FieldViolation {
	var violations []*errdetails.BadRequest_FieldViolation
	for _, id := range ids {
		if err := ValidateUniversalId(id); err != nil {
			violation := &errdetails.BadRequest_FieldViolation{
				Field:       fieldName,
				Description: err.Error(),
			}
			violations = append(violations, violation)
		}
	}
	return violations
}
