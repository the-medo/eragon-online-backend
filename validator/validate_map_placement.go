package validator

import (
	"fmt"
	"github.com/the-medo/talebound-backend/pb"
)

func ValidateMapPlacement(placement *pb.Placement) error {
	inputCount := 0

	if placement.WorldId != nil {
		inputCount++
		err := ValidateWorldId(placement.GetWorldId())
		if err != nil {
			return err
		}
	}

	if placement.QuestId != nil {
		inputCount++
		err := ValidateUniversalId(placement.GetQuestId())
		if err != nil {
			return err
		}
	}

	if inputCount != 1 {
		return fmt.Errorf("exactly one of world_id or quest_id must be provided")
	}

	return nil
}
