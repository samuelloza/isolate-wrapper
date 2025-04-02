package fileSystem

import (
	"fmt"
	"os/exec"
)

func FixOwnership(boxID int) error {
	return nil
	dir := fmt.Sprintf("/var/local/lib/isolate/%d/", boxID)
	cmd := exec.Command("sudo", "chown", "-R", fmt.Sprintf("%s:%s", "sam", "sam"), dir)
	return cmd.Run()
}
