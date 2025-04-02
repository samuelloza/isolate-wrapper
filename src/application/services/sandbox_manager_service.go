package services

import (
	"fmt"
	"os"
	"time"
)

type SandboxManagerService struct {
	MaxBoxes int
}

func NewSandboxManagerService(maxBoxes int) *SandboxManagerService {
	return &SandboxManagerService{MaxBoxes: maxBoxes}
}

func (sm *SandboxManagerService) IsSandboxIDFree(boxID int) bool {
	path := fmt.Sprintf("/var/local/lib/isolate/%d", boxID)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true
	}
	if err != nil {
		return false
	}
	return false
}

func (sm *SandboxManagerService) GetAvailableSandboxID(initialBoxID int) (int, error) {
	boxID := initialBoxID
	attempts := 0

	for {
		isFree := sm.IsSandboxIDFree(boxID)
		if isFree {
			return boxID, nil
		}

		boxID = (boxID + 1) % sm.MaxBoxes
		attempts++

		if attempts >= sm.MaxBoxes {
			return 0, fmt.Errorf("no available sandbox IDs after %d attempts", attempts)
		}

		// Waiting 3 seconds
		time.Sleep(3 * time.Second)
	}
}
