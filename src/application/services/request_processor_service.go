package services

import (
	"fmt"
	"log"

	"github.com/samuelloza/isolate-wrapper/src/application/abstractions"
	"github.com/samuelloza/isolate-wrapper/src/domain"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/comparator"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/compiler"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/fileSystem"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/isolate"
)

type RequestProcessor struct {
	Input            domain.EvaluationInput
	BoxPool          *BoxPool
	TestCaseProvider abstractions.TestCaseProvider
}

func (rp *RequestProcessor) ProcessRequest() (domain.EvaluationResult, error) {
	cases, err := rp.TestCaseProvider.GetTestCases(fmt.Sprintf("%d", rp.Input.ProblemID))
	rp.Input.TestCases = cases

	if err != nil {
		log.Fatalf("Error getting testcases: %v", err)
	}

	SandboxManagerService := NewSandboxManagerService()
	boxId, err := SandboxManagerService.GetAvailableSandboxID(rp.Input.BoxID, rp.BoxPool)
	if err != nil {
		log.Fatalf("Error getting sandbox: %v", err)
	}
	rp.Input.BoxID = boxId
	defer rp.BoxPool.Release(rp.Input.BoxID)

	sandboxImpl := &isolate.IsolateSandbox{}
	compilerImpl, err := compiler.GetCompiler(rp.Input.Language, rp.Input.BoxID)
	if err != nil {
		log.Fatalf("Compiler error: %v", err)
	}
	fileSystemImpl := &fileSystem.FileSystem{}
	comparatorImpl := &comparator.Comparator{}

	evaluator := NewEvaluatorService(sandboxImpl, compilerImpl, fileSystemImpl, comparatorImpl)

	result, err := evaluator.Evaluate(rp.Input)
	if err != nil {
		log.Fatalf("Error evaluating: %v", err)
	}

	//spew.Dump(result)
	return result, nil
}
