package util

import (
	"fmt"
	"github.com/the-medo/talebound-backend/pb"
)

func StringToEvaluationType(value string) (pb.EvaluationType, error) {
	if enumValue, ok := pb.EvaluationType_value[value]; ok {
		return pb.EvaluationType(enumValue), nil
	}
	return pb.EvaluationType_self, fmt.Errorf("invalid EvaluationType: %s", value)
}
