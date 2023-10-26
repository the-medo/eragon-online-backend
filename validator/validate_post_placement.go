package validator

import (
	"fmt"
	"github.com/the-medo/talebound-backend/pb"
)

func ValidatePostPlacement(placement *pb.Placement) error {
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

	if placement.SystemId != nil {
		inputCount++
		err := ValidateUniversalId(placement.GetSystemId())
		if err != nil {
			return err
		}
	}

	if placement.CharacterId != nil {
		inputCount++
		err := ValidateUniversalId(placement.GetCharacterId())
		if err != nil {
			return err
		}
	}

	if inputCount != 1 {
		return fmt.Errorf("exactly one of world_id or quest_id must be provided")
	}

	return nil
}

func ValidatePostPlacementExtended(worldId *int32, questId *int32, systemId *int32, characterId *int32) error {
	inputCount := 0

	if worldId != nil {
		inputCount++
		err := ValidateWorldId(*worldId)
		if err != nil {
			return err
		}
	}

	if questId != nil {
		inputCount++
		err := ValidateUniversalId(*questId)
		if err != nil {
			return err
		}
	}

	if systemId != nil {
		inputCount++
		err := ValidateUniversalId(*systemId)
		if err != nil {
			return err
		}
	}

	if characterId != nil {
		inputCount++
		err := ValidateUniversalId(*characterId)
		if err != nil {
			return err
		}
	}

	if inputCount != 1 {
		return fmt.Errorf("exactly one of the IDs must be provided")
	}

	return nil
}
