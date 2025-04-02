package application

import (
	"fmt"

	"github.com/samuelloza/isolate-wrapper/src/internal/domain/model"
	"github.com/samuelloza/isolate-wrapper/src/internal/infrastructure/comparator"
	"github.com/samuelloza/isolate-wrapper/src/internal/infrastructure/compiler"
	"github.com/samuelloza/isolate-wrapper/src/internal/infrastructure/fileSystem"
	"github.com/samuelloza/isolate-wrapper/src/internal/infrastructure/isolate"
)

type EvaluatorService struct{}

func NewEvaluatorService() *EvaluatorService {
	return &EvaluatorService{}
}

func (s *EvaluatorService) Evaluate(input model.EvaluationInput) (model.EvaluationResult, error) {
	//Default configurations
	sandbox := isolate.IsolateSandbox{
		BoxID:            input.BoxID,
		RunLogPrefix:     input.MetaPrefix,
		FSize:            1048576,
		StackSpace:       65536,
		AddressSpace:     input.RunLimits.Memory,
		StdinFile:        "input.txt",
		StdoutFile:       "user_output.txt",
		StderrFile:       "error.txt",
		Timeout:          float64(input.RunLimits.Time),
		WallclockTimeout: 5.0,
		ExtraTimeout:     0.5,
	}

	//Init Isolate
	_, err := sandbox.Init()
	if err != nil {
		return model.EvaluationResult{}, err
	}
	defer sandbox.Cleanup()

	//Get name and ext to compile the code
	srcName := fmt.Sprintf("Main.%s", input.Language)

	//Create a tmp directory "/tmp/patito-wrapper-BoxID"
	fileSystem.CreateTmpDirectory(input.BoxID)

	//Copy the code in tmp directory
	_ = fileSystem.WriteCodeInTmpDirectory(input.SourceCode, input.BoxID, srcName)

	//Compile code
	outputDir := fmt.Sprintf("/tmp/patito-wrapper-%d", input.BoxID)
	fmt.Printf("Tmp directory %s\n", outputDir)

	comp, err := compiler.GetCompiler(input.Language, outputDir)
	if err != nil {
		return model.EvaluationResult{}, fmt.Errorf("compilation error: %w", err)
	}

	err = comp.Compile("Main."+input.Language, outputDir)
	if err != nil {
		return model.EvaluationResult{}, fmt.Errorf("compilation failed: %w", err)
	}

	var results []model.TestCaseResult
	totalPassed := 0
	for i, test := range input.TestCases {
		_ = fileSystem.CopyFileToTmpDirectory(test.Input, input.BoxID, "input.txt")
		_ = fileSystem.CopyFileToTmpDirectory(test.Output, input.BoxID, "expected.txt")

		logData, err := sandbox.Run(fmt.Sprintf("%s/Main", outputDir), i)

		outputPath := fmt.Sprintf("/var/local/lib/isolate/%d/box/user_output.txt", input.BoxID)
		expectedPath := fmt.Sprintf("/tmp/patito-wrapper-%d/expected.txt", input.BoxID)
		inputPath := fmt.Sprintf("/tmp/patito-wrapper-%d/input.txt", input.BoxID)
		// errorPath := fmt.Sprintf("/var/local/lib/isolate/%d/box/error.txt", input.BoxID)

		cmpResult, _ := comparator.CompareZOJ(expectedPath, outputPath, inputPath)

		passed := cmpResult == comparator.OJ_AC || cmpResult == comparator.OJ_PE
		received := outputPath

		execTime := 0
		memUsed := 0

		if val, ok := logData["time"]; ok {
			fmt.Sscanf(val, "%d", &execTime)
		}
		if val, ok := logData["max-rss"]; ok {
			fmt.Sscanf(val, "%d", &memUsed)
		}

		res := model.TestCaseResult{
			Index:         i,
			Passed:        passed,
			Expected:      test.Output,
			Received:      received,
			ExecutionTime: execTime,
			MemoryUsed:    memUsed,
		}
		if err != nil {
			res.ErrorMessage = err.Error()
		}
		if passed {
			totalPassed++
		}
		results = append(results, res)
	}

	status := fmt.Sprintf("%d/%d", totalPassed, len(results))
	if totalPassed == len(results) {
		status = "Accepted"
	}

	return model.EvaluationResult{
		EvaluationID: input.ID,
		Results:      results,
		TotalPassed:  totalPassed,
		TotalCases:   len(results),
		Status:       status,
	}, nil
}
