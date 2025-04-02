package isolate

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/samuelloza/isolate-wrapper/src/application/abstractions"
)

type IsolateSandbox struct{}

func (s *IsolateSandbox) Init(boxID int) error {
	cmd := exec.Command("/usr/local/bin/isolate", "--box-id", fmt.Sprint(boxID), "--cg", "--init")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to initialize sandbox: %s", string(output))
	}
	return nil
}

func (s *IsolateSandbox) Cleanup(boxID int) error {
	cmd := exec.Command("/usr/local/bin/isolate", "--box-id", fmt.Sprint(boxID), "--cg", "--cleanup")
	_, err := cmd.Output()
	return err
}

func (s *IsolateSandbox) Run(boxID int, caseIndex int) (abstractions.SandboxLogData, error) {
	opts := s.BuildBoxOptions(boxID, caseIndex)
	opts = append(opts, "--", "/source/Main")

	cmd := exec.Command("/usr/local/bin/isolate", opts...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", output)
		fmt.Printf("Error to execute Box id %d:\n%s\n", boxID, string(err.Error())+string(output))
		return abstractions.SandboxLogData{}, err
	}

	log, err := s.ReadLog(boxID)
	if err != nil {
		return abstractions.SandboxLogData{}, err
	}

	execTime, _ := strconv.Atoi(log["time"])
	memUsed, _ := strconv.Atoi(log["cg-mem"])

	return abstractions.SandboxLogData{
		ExecutionTime: execTime,
		MemoryUsed:    memUsed,
	}, nil
}

func (s *IsolateSandbox) BuildBoxOptions(boxID int, caseIndex int) []string {
	tmpDirectory := fmt.Sprintf("/tmp/patito-wrapper-%d", boxID)
	sandboxPath := "/source"
	return []string{
		fmt.Sprintf("--box-id=%d", boxID),
		"--cg",
		fmt.Sprintf("--dir=/source=%s:rw", tmpDirectory),
		fmt.Sprintf("--stdin=%s", fmt.Sprintf("%s/%s", sandboxPath, "input.txt")),
		fmt.Sprintf("--stdout=%s", fmt.Sprintf("/box/%s", "user_output.txt")),
		fmt.Sprintf("--stderr=%s", fmt.Sprintf("/box/%s", "error.txt")),
		fmt.Sprintf("--meta=%s", fmt.Sprintf("%s/meta", tmpDirectory)),
		fmt.Sprintf("--fsize=%d", 1024),
		fmt.Sprintf("--cg-mem=%d", 1024),
		fmt.Sprintf("--time=%.2f", 1.0),
		fmt.Sprintf("--wall-time=%.2f", 1.0),
		fmt.Sprintf("--extra-time=%.2f", 0.5),
		"--run",
	}
}

func (s *IsolateSandbox) ReadLog(boxID int) (map[string]string, error) {
	logFile := fmt.Sprintf("/tmp/patito-wrapper-%d/meta", boxID)
	data, err := os.ReadFile(logFile)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result, nil
}
