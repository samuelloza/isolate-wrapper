package services

import (
	"fmt"
	"os"
	"time"
)

type SandboxManagerService struct{}

func NewSandboxManagerService() *SandboxManagerService {
	return &SandboxManagerService{}
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

func (sm *SandboxManagerService) GetAvailableSandboxID(initialBoxID int, boxPool *BoxPool) (int, error) {
	maxAttempts := 20

	for i := 0; i < maxAttempts; i++ {
		select {
		case boxID := <-boxPool.pool:
			return boxID, nil
		default:
			time.Sleep(3 * time.Second)
		}
	}

	return 0, fmt.Errorf("timeout: all sandboxes are busy after %d attempts", maxAttempts)
}
