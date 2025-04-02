package application

import (
	"os"
	"testing"
	"time"

	"github.com/samuelloza/isolate-wrapper/src/application/services"
)

func Test_SandboxManagerService_GetAvailableSandboxID(t *testing.T) {
	// Setup: Create a mock sandbox directory to simulate an unavailable sandbox
	unavailableBoxID := 1
	mockPath := "/var/local/lib/isolate/1"
	if err := os.MkdirAll(mockPath, 0755); err != nil {
		t.Fatalf("Failed to create mock sandbox directory: %v", err)
	}
	defer os.RemoveAll(mockPath) // Cleanup after test

	// Initialize SandboxManagerService with a small number of boxes
	manager := services.NewSandboxManagerService(3)

	// Test: Request an available sandbox ID starting from the unavailable one
	boxID, err := manager.GetAvailableSandboxID(unavailableBoxID)
	if err != nil {
		t.Fatalf("Failed to get available sandbox ID: %v", err)
	}

	// Verify: Ensure the returned box ID is different from the unavailable one
	if boxID == unavailableBoxID {
		t.Errorf("Expected a different sandbox ID, but got the same: %d", boxID)
	}

	// Verify: Ensure the returned box ID is within the valid range
	if boxID < 0 || boxID >= 3 {
		t.Errorf("Returned sandbox ID %d is out of range", boxID)
	}

	// Simulate waiting for sandbox availability
	time.Sleep(1 * time.Second)
}
