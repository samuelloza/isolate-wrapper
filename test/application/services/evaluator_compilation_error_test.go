package services_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/samuelloza/isolate-wrapper/src/application/abstractions"
	"github.com/samuelloza/isolate-wrapper/src/application/services"
	"github.com/samuelloza/isolate-wrapper/src/domain"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/comparator"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/compiler"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/fileSystem"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/isolate"
)

const sumSourcece = `#include <iostream>
#include <sstream>
#include <string>
using namespace std;
int main(){
    string line;
    if(line=="RTE"){
         int a = 0;
}
`

func TestEvaluator_CE(t *testing.T) {
	t.Parallel()
	cwd, _ := os.Getwd()
	projectRoot := filepath.Join(cwd, "../../../")

	paths := []struct{ in, out string }{
		{"test/testcases/1000/sum_allac_1.in", "test/testcases/1000/sum_allac_1.out"},
	}
	testCases := getTestCases(paths, projectRoot)

	input := domain.EvaluationInput{
		ID:         "all-ac",
		UniqID:     "all-ac",
		BoxID:      0,
		ProblemID:  1000,
		Language:   "cpp",
		SourceCode: sumSourcece,
		RunLimits:  domain.RunLimits{Time: 2, Memory: 65536, Output: 1024},
		TestCases:  testCases,
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

	compilerService, _ := compiler.GetCompiler(input.Language, input.BoxID)
	if err != nil {
		t.Logf("Error getting compiler: %v", err)
	}
	evaluator := services.NewEvaluatorService(sandbox, compilerService, fs, cmp)

	result, err := evaluator.Evaluate(input)
	defer boxPool.Release(input.BoxID)

	if err != nil {
		t.Logf("Error during evaluation: %v", err)
	}

	spew.Dump(result)

	if result.TotalPassed != 0 {
		t.Errorf("Expected %d test cases to pass, but %d passed", len(result.Results), result.TotalPassed)
	}
	if result.Status != abstractions.OJ_CE {
		t.Errorf("Expected status %d (OJ_CE), but got %d", domain.OJ_CE, result.Status)
	}
}
