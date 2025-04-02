package application

import (
	"fmt"

	"github.com/samuelloza/isolate-wrapper/src/internal/domain/interfaces"
	"github.com/samuelloza/isolate-wrapper/src/internal/domain/model"
)

type EvaluatorService struct {
	Sandbox    interfaces.Sandbox
	Compiler   interfaces.Compiler
	FileSystem interfaces.FileSystem
	Comparator interfaces.Comparator
}

func NewEvaluatorService(sandbox interfaces.Sandbox, compiler interfaces.Compiler, fileSystem interfaces.FileSystem, comparator interfaces.Comparator) *EvaluatorService {
	return &EvaluatorService{
		Sandbox:    sandbox,
		Compiler:   compiler,
		FileSystem: fileSystem,
		Comparator: comparator,
	}
}

func (s *EvaluatorService) Evaluate(input model.EvaluationInput) (model.EvaluationResult, error) {
	// Initialize sandbox
	if err := s.Sandbox.Init(input.BoxID); err != nil {
		return model.EvaluationResult{}, fmt.Errorf("failed to initialize sandbox: %w", err)
	}
	defer s.Sandbox.Cleanup(input.BoxID)

	// Write source code to temporary directory
	srcFileName := fmt.Sprintf("Main.%s", input.Language)
	if err := s.FileSystem.WriteFile(input.BoxID, srcFileName, input.SourceCode); err != nil {
		return model.EvaluationResult{}, fmt.Errorf("failed to write source code: %w", err)
	}

	// Compile source code
	outputDir := fmt.Sprintf("/tmp/patito-wrapper-%d", input.BoxID)
	if err := s.Compiler.Compile(srcFileName, outputDir); err != nil {
		return model.EvaluationResult{}, fmt.Errorf("compilation failed: %w", err)
	}

	// Execute test cases
	var results []model.TestCaseResult
	totalPassed := 0

	for i, test := range input.TestCases {
		// Prepare input and expected output files
		if err := s.FileSystem.CopyFile(test.Input, input.BoxID, "input.txt"); err != nil {
			return model.EvaluationResult{}, fmt.Errorf("failed to copy input file: %w", err)
		}

		if err := s.FileSystem.CopyFile(test.Output, input.BoxID, "expected.txt"); err != nil {
			return model.EvaluationResult{}, fmt.Errorf("failed to copy expected output file: %w", err)
		}

		// Run the program in the sandbox
		logData, err := s.Sandbox.Run(input.BoxID, i)
		if err != nil {
			results = append(results, model.TestCaseResult{
				Index:        i,
				Passed:       false,
				ErrorMessage: err.Error(),
			})
			continue
		}

		// Compare outputs
		outputPath := s.FileSystem.GetOutputPath(input.BoxID)
		expectedPath := s.FileSystem.GetFilePath(input.BoxID, "expected.txt")
		cmpResult, err := s.Comparator.Compare(expectedPath, outputPath)
		if err != nil {
			results = append(results, model.TestCaseResult{
				Index:        i,
				Passed:       false,
				ErrorMessage: err.Error(),
			})
			continue
		}

		// Record result
		passed := cmpResult == model.OJ_AC
		if passed {
			totalPassed++
		}
		results = append(results, model.TestCaseResult{
			Index:         i,
			Passed:        passed,
			Expected:      test.Output,
			Received:      outputPath,
			ExecutionTime: logData.ExecutionTime,
			MemoryUsed:    logData.MemoryUsed,
		})
	}

	// Determine overall status
	status := "Accepted"
	if totalPassed != len(results) {
		status = fmt.Sprintf("%d/%d", totalPassed, len(results))
	}

	return model.EvaluationResult{
		EvaluationID: input.ID,
		Results:      results,
		TotalPassed:  totalPassed,
		TotalCases:   len(results),
		Status:       status,
	}, nil
}
