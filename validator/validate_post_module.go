package validator

import (
	"fmt"
	"github.com/the-medo/talebound-backend/pb"
)

func ValidatePostModule(module *pb.Module) error {
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

	if module.SystemId != nil {
		inputCount++
		err := ValidateUniversalId(module.GetSystemId())
		if err != nil {
			return err
		}
	}

	if module.CharacterId != nil {
		inputCount++
		err := ValidateUniversalId(module.GetCharacterId())
		if err != nil {
			return err
		}
	}

	if inputCount != 1 {
		return fmt.Errorf("exactly one of world_id or quest_id must be provided")
	}

	return nil
}

func ValidatePostModuleExtended(worldId *int32, questId *int32, systemId *int32, characterId *int32) error {
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
