//revive:disable:var-naming For now it's okay to have a generic name
package utils

import (
	"fmt"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/types"
)

func DetectedBlockTypeFromHcl(hclType string) (*types.DetectedBlockType, error) {
	switch hclType {
	case "resource":
		return &constants.DetectedBlockTypeResource, nil
	case "data":
		return &constants.DetectedBlockTypeData, nil
	case "cloud":
		return &constants.DetectedBlockTypeCloud, nil
	case "backend":
		return &constants.DetectedBlockTypeBackend, nil
	default:
		return nil, fmt.Errorf("unknown detected block type: %s", hclType)
	}
}
