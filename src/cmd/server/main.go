package main

import (
	"log"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/samuelloza/isolate-wrapper/src/application/services"
	"github.com/samuelloza/isolate-wrapper/src/domain"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/comparator"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/compiler"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/fileSystem"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/isolate"
)

func main() {
	cwd := "/home/sam/project/github/isolate-wrapper"

	input := domain.EvaluationInput{
		ID:          "123456",
		UniqID:      "abc123",
		BoxID:       0,
		ProblemName: "Suma de Dos Números",
		Language:    "cpp",
		SourceCode: `#include <iostream>
					using namespace std;
					int main() {
					    int a;
					    cin>>a;
					    cout<<a*10<<endl;
					    return 0;
					}`,
		RunLimits: domain.RunLimits{
			Time:   2,
			Memory: 65536,
			Output: 1024,
		},
		TestCases: []domain.TestCase{
			{Input: filepath.Join(cwd, "test/testcases/test_1.in"), Output: filepath.Join(cwd, "test/testcases/test_1.out")},
			{Input: filepath.Join(cwd, "test/testcases/test_2.in"), Output: filepath.Join(cwd, "test/testcases/test_2.out")},
		},
	}

	boxPool := services.NewBoxPool(2)

	SandboxManagerService := services.NewSandboxManagerService()
	boxId, err := SandboxManagerService.GetAvailableSandboxID(input.BoxID, boxPool)
	if err != nil {
		log.Fatalf("Error getting sandbox: %v", err)
	}
	input.BoxID = boxId

	sandboxImpl := &isolate.IsolateSandbox{}
	compilerImpl, err := compiler.GetCompiler(input.Language, input.BoxID)
	if err != nil {
		log.Fatalf("Compiler error: %v", err)
	}
	fileSystemImpl := &fileSystem.FileSystem{}
	comparatorImpl := &comparator.Comparator{}

	evaluator := services.NewEvaluatorService(sandboxImpl, compilerImpl, fileSystemImpl, comparatorImpl)

	result, err := evaluator.Evaluate(input)
	if err != nil {
		log.Fatalf("Error evaluating: %v", err)
	}

	spew.Dump(result)
}
