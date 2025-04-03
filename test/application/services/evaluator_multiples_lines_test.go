package services_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/samuelloza/isolate-wrapper/src/application/services"
	"github.com/samuelloza/isolate-wrapper/src/domain"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/comparator"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/compiler"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/fileSystem"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/isolate"
)

func TestEvaluator_MultipleLines(t *testing.T) {
	t.Parallel()
	cwd, _ := os.Getwd()
	projectRoot := filepath.Join(cwd, "../../../")

	paths := []struct{ in, out string }{
		{"test/testcases/multi_case_1.in", "test/testcases/multi_case_1.out"},
	}
	testCases := getTestCases(paths, projectRoot)

	input := domain.EvaluationInput{
		ID:          "multi-lines",
		UniqID:      "multi-lines",
		BoxID:       0,
		ProblemName: "Multiple lines",
		Language:    "cpp",
		SourceCode: `#include <iostream>
using namespace std;
int main() {
    int x;
    while (cin >> x) {
        cout << x * 2 << endl;
    }
    return 0;
}
`,
		RunLimits: domain.RunLimits{Time: 2, Memory: 65536, Output: 2048},
		TestCases: testCases,
	}

	sandbox := &isolate.IsolateSandbox{}
	fs := &fileSystem.FileSystem{}
	cmp := &comparator.Comparator{}

	sandboxManager := services.NewSandboxManagerService()
	availableID, err := sandboxManager.GetAvailableSandboxID(input.BoxID, boxPool)
	if err != nil {
		t.Fatalf("No sandbox available: %v", err)
	}
	input.BoxID = availableID

	compilerService, err := compiler.GetCompiler(input.Language, input.BoxID)
	if err != nil {
		t.Logf("Error getting compiler: %v", err)
	}
	evaluator := services.NewEvaluatorService(sandbox, compilerService, fs, cmp)

	result, err := evaluator.Evaluate(input)
	defer boxPool.Release(input.BoxID)

	if err != nil {
		t.Fatalf("Error to evaluate: %v", err)
	}

	spew.Dump(result)

	if result.TotalPassed != len(result.Results) {
		t.Errorf("Expected all to pass, but passed %d of %d", result.TotalPassed, len(result.Results))
	}
}

func TestEvaluator_OutputLimitExceeded(t *testing.T) {
	t.Parallel()
	cwd, _ := os.Getwd()
	projectRoot := filepath.Join(cwd, "../../../")

	paths := []struct{ in, out string }{
		{"test/testcases/output_limit.in", "test/testcases/output_limit.out"},
	}
	testCases := getTestCases(paths, projectRoot)

	input := domain.EvaluationInput{
		ID:          "output-limit",
		UniqID:      "output-limit",
		BoxID:       0,
		ProblemName: "Output too large",
		Language:    "cpp",
		SourceCode: `#include <iostream>
using namespace std;
int main() {
    for (int i = 0; i < 1000000; ++i) {
        cout << "SPAM" << endl;
    }
    return 0;
}
`,
		RunLimits: domain.RunLimits{Time: 2, Memory: 65536, Output: 128},
		TestCases: testCases,
	}

	sandbox := &isolate.IsolateSandbox{}
	fs := &fileSystem.FileSystem{}
	cmp := &comparator.Comparator{}

	sandboxManager := services.NewSandboxManagerService()
	availableID, err := sandboxManager.GetAvailableSandboxID(input.BoxID, boxPool)
	if err != nil {
		t.Fatalf("No sandbox available: %v", err)
	}
	input.BoxID = availableID

	compilerService, err := compiler.GetCompiler(input.Language, input.BoxID)
	if err != nil {
		t.Fatalf("Error getting compiler: %v", err)
	}
	evaluator := services.NewEvaluatorService(sandbox, compilerService, fs, cmp)

	result, err := evaluator.Evaluate(input)
	defer boxPool.Release(input.BoxID)

	if err != nil {
		t.Fatalf("Error to evaluate: %v", err)
	}

	spew.Dump(result)

	if result.Results[0].Status != domain.OJ_OL {
		t.Errorf("The response should be OJ_OL, but got %d", result.Status)
	}
}
