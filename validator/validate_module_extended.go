package validator

import (
	"fmt"
)

func ValidateModuleExtended(worldId *int32, questId *int32, systemId *int32, characterId *int32) error {
	inputCount := 0

	if worldId != nil && *worldId > 0 {
		inputCount++
		err := ValidateModuleId(*worldId)
		if err != nil {
			return err
		}
	}

	if questId != nil && *questId > 0 {
		inputCount++
		err := ValidateUniversalId(*questId)
		if err != nil {
			return err
		}
	}

	if systemId != nil && *systemId > 0 {
		inputCount++
		err := ValidateUniversalId(*systemId)
		if err != nil {
			return err
		}
	}

	if characterId != nil && *characterId > 0 {
		inputCount++
		err := ValidateUniversalId(*characterId)
		if err != nil {
			return err
		}
	}

	if inputCount != 1 {
		return fmt.Errorf("exactly one of the IDs must be provided, got %d", inputCount)
	}

	return nil
}
