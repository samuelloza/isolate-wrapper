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

const sumExtendedSource = `#include <iostream>
#include <sstream>
#include <string>
using namespace std;
int main(){
    string line;
    getline(cin, line);

	if(line=="WRONG"){
         return 0;
    }
    if(line=="RTE"){
         cout << 1/0;
         return 0;
    }
    if(line=="TLE"){
         while(true){}
         return 0;
    }
    istringstream iss(line);
    int a, b;
    if(!(iss >> a >> b)){
         cout << "WRONG";
         return 0;
    }
    cout << (a + b) << endl;
    return 0;
}
`

func TestEvaluator_OnlyExtendedCases(t *testing.T) {
	t.Parallel()
	cwd, _ := os.Getwd()
	projectRoot := filepath.Join(cwd, "../../../")

	paths := []struct{ in, out string }{
		{"test/testcases/extended_case_1.in", "test/testcases/extended_case_1.out"},
		{"test/testcases/extended_case_2.in", "test/testcases/extended_case_2.out"},
		{"test/testcases/extended_case_3.in", "test/testcases/extended_case_3.out"},
		{"test/testcases/extended_case_4.in", "test/testcases/extended_case_4.out"},
		{"test/testcases/extended_case_5.in", "test/testcases/extended_case_5.out"},
		{"test/testcases/extended_case_6.in", "test/testcases/extended_case_6.out"},
		{"test/testcases/extended_case_7.in", "test/testcases/extended_case_7.out"},
		{"test/testcases/extended_case_8.in", "test/testcases/extended_case_8.out"},
	}

	solutions := []int{
		abstractions.OJ_AC,
		abstractions.OJ_AC,
		abstractions.OJ_WA,
		abstractions.OJ_RE,
		abstractions.OJ_TL,
		abstractions.OJ_AC,
		abstractions.OJ_AC,
		abstractions.OJ_AC,
	}

	testCases := getTestCases(paths, projectRoot)

	input := domain.EvaluationInput{
		ID:          "2",
		UniqID:      "2",
		BoxID:       0,
		ProblemName: "Extended Test",
		Language:    "cpp",
		SourceCode:  sumExtendedSource,
		RunLimits:   domain.RunLimits{Time: 1, Memory: 65536, Output: 1024},
		TestCases:   testCases,
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

	for i, result := range result.Results {
		if result.Status != solutions[i] {
			t.Errorf("Test %d: Expected status %d, but got %d", i+1, solutions[i], result.Status)
		}
	}

	if result.Status != abstractions.OJ_WA {
		t.Errorf("Expected overall status 6 (OJ_WA), but got %d", result.Status)
	}
}
