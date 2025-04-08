package testcaseprovider

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/samuelloza/isolate-wrapper/src/application/abstractions"
	"github.com/samuelloza/isolate-wrapper/src/domain"
)

type FileSystemTestCaseProvider struct {
	BasePath string
}

func NewFileSystemTestCaseProvider(basePath string) abstractions.TestCaseProvider {
	return &FileSystemTestCaseProvider{BasePath: basePath}
}

func (p *FileSystemTestCaseProvider) GetTestCases(problemID string) ([]domain.TestCase, error) {
	testCases := []domain.TestCase{}
	dir := filepath.Join(p.BasePath, problemID)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read test case directory: %w", err)
	}

	inputPaths := make(map[string]string)
	outputPaths := make(map[string]string)

	for _, file := range files {
		name := file.Name()
		fullPath := filepath.Join(dir, name)

		if strings.HasSuffix(name, ".in") {
			key := strings.ReplaceAll(name, ".in", "")
			inputPaths[key] = fullPath
		} else if strings.HasSuffix(name, ".out") {
			key := strings.ReplaceAll(name, ".out", "")
			outputPaths[key] = fullPath
		}
	}

	keys := []string{}
	for k := range inputPaths {
		if _, ok := outputPaths[k]; ok {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		tc := domain.TestCase{
			Input:  inputPaths[k],
			Output: outputPaths[k],
		}
		testCases = append(testCases, tc)
	}

	return testCases, nil
}
