package services_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/samuelloza/isolate-wrapper/src/application/abstractions"
	"github.com/samuelloza/isolate-wrapper/src/application/services"
	"github.com/samuelloza/isolate-wrapper/src/domain"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/comparator"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/compiler"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/fileSystem"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/isolate"
)

const sumSource = `#include <iostream>
#include <sstream>
#include <string>
using namespace std;
int main(){
    string line;
    getline(cin, line);
    if(line=="RTE"){
         int a = 0;
         cout << 1/a;
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

func TestEvaluator_AllAC(t *testing.T) {
	t.Parallel()
	cwd, _ := os.Getwd()
	projectRoot := filepath.Join(cwd, "../../../")

	paths := []struct{ in, out string }{
		{"test/testcases/sum_allac_1.in", "test/testcases/sum_allac_1.out"},
		{"test/testcases/sum_allac_2.in", "test/testcases/sum_allac_2.out"},
	}
	testCases := getTestCases(paths, projectRoot)

	input := domain.EvaluationInput{
		ID:          "all-ac",
		UniqID:      "all-ac",
		BoxID:       0,
		ProblemName: "A + B",
		Language:    "cpp",
		SourceCode:  sumSource,
		RunLimits:   domain.RunLimits{Time: 2, Memory: 65536, Output: 1024},
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
		t.Fatalf("Error during evaluation: %v", err)
	}

	spew.Dump(result)

	if result.TotalPassed != len(result.Results) {
		t.Errorf("Expected %d test cases to pass, but %d passed", len(result.Results), result.TotalPassed)
	}
	if result.Status != abstractions.OJ_AC {
		t.Errorf("Expected status 7 (OJ_AC), but got %d", result.Status)
	}
}

func TestEvaluator_Mixed2AC2WA(t *testing.T) {
	t.Parallel()
	cwd, _ := os.Getwd()
	projectRoot := filepath.Join(cwd, "../../../")
	paths := []struct{ in, out string }{
		{"test/testcases/sum_mixed_1.in", "test/testcases/sum_mixed_1.out"},
		{"test/testcases/sum_mixed_2.in", "test/testcases/sum_mixed_2.out"},
		{"test/testcases/sum_mixed_3.in", "test/testcases/sum_mixed_3.out"},
		{"test/testcases/sum_mixed_4.in", "test/testcases/sum_mixed_4.out"},
	}
	testCases := getTestCases(paths, projectRoot)

	input := domain.EvaluationInput{
		ID:          "mixed-2ac-2wa",
		UniqID:      "mixed-2ac-2wa",
		BoxID:       0,
		ProblemName: "A + B",
		Language:    "cpp",
		SourceCode:  sumSource,
		RunLimits:   domain.RunLimits{Time: 2, Memory: 65536, Output: 1024},
		TestCases:   testCases,
	}

	sandboxManager := services.NewSandboxManagerService()
	availableID, err := sandboxManager.GetAvailableSandboxID(input.BoxID, boxPool)
	if err != nil {
		t.Fatalf("No sandbox available: %v", err)
	}
	input.BoxID = availableID

	sandbox := &isolate.IsolateSandbox{}
	compilerService, err := compiler.GetCompiler(input.Language, input.BoxID)
	if err != nil {
		t.Fatalf("Error getting compiler: %v", err)
	}
	fs := &fileSystem.FileSystem{}
	cmp := &comparator.Comparator{}
	evaluator := services.NewEvaluatorService(sandbox, compilerService, fs, cmp)

	result, err := evaluator.Evaluate(input)
	defer boxPool.Release(input.BoxID)

	if err != nil {
		t.Fatalf("Error during evaluation: %v", err)
	}
	spew.Dump(result)

	if result.TotalPassed != 2 && result.TotalCases == 4 {
		t.Errorf("Expected %d test cases to pass, but %d passed", len(result.Results), result.TotalPassed)
	}
	if result.TotalPassed == len(result.Results) {
		t.Errorf("Expected status 7 (OJ_AC), but got %d", result.Status)
	}
}

func TestEvaluator_MixedAC_RTE_TLE(t *testing.T) {
	t.Parallel()
	cwd, _ := os.Getwd()
	projectRoot := filepath.Join(cwd, "../../../")

	paths := []struct{ in, out string }{
		{"test/testcases/sum_mixed_ac_rte_tle_1.in", "test/testcases/sum_mixed_ac_rte_tle_1.out"},
		{"test/testcases/sum_mixed_ac_rte_tle_2.in", "test/testcases/sum_mixed_ac_rte_tle_2.out"},
		{"test/testcases/sum_mixed_ac_rte_tle_3.in", "test/testcases/sum_mixed_ac_rte_tle_3.out"},
	}
	testCases := getTestCases(paths, projectRoot)

	input := domain.EvaluationInput{
		ID:          "mixed-ac-rte-tle",
		UniqID:      "mixed-ac-rte-tle",
		BoxID:       0,
		ProblemName: "A + B",
		Language:    "cpp",
		SourceCode:  sumSource,
		RunLimits:   domain.RunLimits{Time: 2, Memory: 65536, Output: 1024},
		TestCases:   testCases,
	}

	sandboxManager := services.NewSandboxManagerService()
	availableID, err := sandboxManager.GetAvailableSandboxID(input.BoxID, boxPool)
	if err != nil {
		t.Fatalf("No sandbox available: %v", err)
	}
	input.BoxID = availableID

	sandbox := &isolate.IsolateSandbox{}
	compilerService, err := compiler.GetCompiler(input.Language, input.BoxID)
	if err != nil {
		t.Fatalf("Error getting compiler: %v", err)
	}
	fs := &fileSystem.FileSystem{}
	cmp := &comparator.Comparator{}
	evaluator := services.NewEvaluatorService(sandbox, compilerService, fs, cmp)

	done := make(chan struct{})
	var result domain.EvaluationResult
	go func() {
		result, err = evaluator.Evaluate(input)
		defer boxPool.Release(input.BoxID)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatalf("Timeout during execution, suspected TLE in the test")
	}

	if err != nil {
		t.Fatalf("Error during evaluation: %v", err)
	}
	spew.Dump(result)

	if result.TotalPassed != 1 {
		t.Errorf("Expected %d test cases to pass, but %d passed", len(result.Results), result.TotalPassed)
	}
	if result.Status == 7 {
		t.Errorf("Expected status 7 (OJ_AC), but got %d", result.Status)
	}
}

func TestEvaluator_Complex(t *testing.T) {
	t.Parallel()
	cwd, _ := os.Getwd()
	projectRoot := filepath.Join(cwd, "../../../")

	paths := []struct{ in, out string }{
		{"test/testcases/sum_complex_1.in", "test/testcases/sum_complex_1.out"},
		{"test/testcases/sum_complex_2.in", "test/testcases/sum_complex_2.out"},
		{"test/testcases/sum_complex_3.in", "test/testcases/sum_complex_3.out"},
		{"test/testcases/sum_complex_4.in", "test/testcases/sum_complex_4.out"},
		{"test/testcases/sum_complex_5.in", "test/testcases/sum_complex_5.out"},
		{"test/testcases/sum_complex_6.in", "test/testcases/sum_complex_6.out"},
		{"test/testcases/sum_complex_7.in", "test/testcases/sum_complex_7.out"},
	}
	testCases := getTestCases(paths, projectRoot)

	input := domain.EvaluationInput{
		ID:          "complex",
		UniqID:      "complex",
		BoxID:       0,
		ProblemName: "A + B",
		Language:    "cpp",
		SourceCode:  sumSource,
		RunLimits:   domain.RunLimits{Time: 2, Memory: 65536, Output: 1024},
		TestCases:   testCases,
	}

	sandboxManager := services.NewSandboxManagerService()
	availableID, err := sandboxManager.GetAvailableSandboxID(input.BoxID, boxPool)
	if err != nil {
		t.Fatalf("No sandbox available: %v", err)
	}
	input.BoxID = availableID

	sandbox := &isolate.IsolateSandbox{}
	compilerService, err := compiler.GetCompiler(input.Language, input.BoxID)
	if err != nil {
		t.Fatalf("Error getting compiler: %v", err)
	}
	fs := &fileSystem.FileSystem{}
	cmp := &comparator.Comparator{}
	evaluator := services.NewEvaluatorService(sandbox, compilerService, fs, cmp)

	done := make(chan struct{})
	var result domain.EvaluationResult
	go func() {
		result, err = evaluator.Evaluate(input)
		defer boxPool.Release(input.BoxID)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(7 * time.Second):
		t.Fatalf("Timeout during execution, suspected TLE in the test")
	}

	if err != nil {
		t.Fatalf("Error in the evaluation: %v", err)
	}
	spew.Dump(result)

	if result.TotalPassed != 1 {
		t.Errorf("Expected 1 test case to pass, but %d passed", result.TotalPassed)
	}
	if result.Status == 7 {
		t.Errorf("Expected status different from 7 (OJ_OE), but got %d", result.Status)
	}
}
