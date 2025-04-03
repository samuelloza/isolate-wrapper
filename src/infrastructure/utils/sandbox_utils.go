package utils

import (
	"fmt"
	"os"
)

func IsSandboxIDFree(boxID int) bool {
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
