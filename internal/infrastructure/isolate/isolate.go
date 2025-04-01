package isolate

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type IsolateSandbox struct {
	BoxID            int
	RunLogPrefix     string
	FSize            int
	StackSpace       int
	AddressSpace     int
	StdinFile        string
	StdoutFile       string
	StderrFile       string
	Timeout          float64
	WallclockTimeout float64
	ExtraTimeout     float64
}

func (s *IsolateSandbox) Init() (string, error) {
	cmd := exec.Command("/usr/local/bin/isolate", "--box-id", fmt.Sprint(s.BoxID), "--cg", "--init")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error to Init isolate Box id %d:\n%s\n", s.BoxID, string(output))
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (s *IsolateSandbox) Cleanup() (string, error) {
	cmd := exec.Command("/usr/local/bin/isolate", "--box-id", fmt.Sprint(s.BoxID), "--cg", "--cleanup")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (s *IsolateSandbox) BuildBoxOptions(caseIndex int) []string {
	tmpDirectory := fmt.Sprintf("/tmp/patito-wrapper-%d", s.BoxID)
	sandboxPath := "/source"
	return []string{
		fmt.Sprintf("--box-id=%d", s.BoxID),
		"--cg",
		fmt.Sprintf("--dir=/source=%s:rw", tmpDirectory),
		fmt.Sprintf("--stdin=%s", fmt.Sprintf("%s/%s", sandboxPath, s.StdinFile)),
		fmt.Sprintf("--stdout=%s", fmt.Sprintf("/box/%s", s.StdoutFile)),
		fmt.Sprintf("--stderr=%s", fmt.Sprintf("/box/%s", s.StderrFile)),
		fmt.Sprintf("--meta=%s", fmt.Sprintf("%s/meta", tmpDirectory)),
		fmt.Sprintf("--fsize=%d", s.FSize/1024),
		fmt.Sprintf("--cg-mem=%d", s.AddressSpace),
		fmt.Sprintf("--time=%.2f", s.Timeout),
		fmt.Sprintf("--wall-time=%.2f", s.WallclockTimeout),
		fmt.Sprintf("--extra-time=%.2f", s.ExtraTimeout),
		"--run",
	}
}

func (s *IsolateSandbox) Run(executable string, caseIndex int) (map[string]string, error) {
	opts := s.BuildBoxOptions(caseIndex)
	opts = append(opts, "--", "/source/Main")

	cmd := exec.Command("/usr/local/bin/isolate", opts...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("%s\n", output)
		fmt.Printf("Error to execute Box id %d:\n%s\n", s.BoxID, string(err.Error())+string(output))
		return nil, err
	}

	return s.ReadLog(s.BoxID)
}

func (s *IsolateSandbox) ReadLog(boxId int) (map[string]string, error) {
	logFile := fmt.Sprintf("/tmp/patito-wrapper-%d/meta", boxId)
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
