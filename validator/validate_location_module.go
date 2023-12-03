package validator

import (
	"fmt"
	"github.com/the-medo/talebound-backend/pb"
)

func ValidateLocationModule(module *pb.ModuleDefinition) error {
	inputCount := 0

	if module.WorldId != nil {
		inputCount++
		err := ValidateWorldId(module.GetWorldId())
		if err != nil {
			return err
		}
	}

	if module.QuestId != nil {
		inputCount++
		err := ValidateUniversalId(module.GetQuestId())
		if err != nil {
			return err
		}
	}

	if inputCount != 1 {
		return fmt.Errorf("exactly one of world_id or quest_id must be provided")
	}

	return nil
}
