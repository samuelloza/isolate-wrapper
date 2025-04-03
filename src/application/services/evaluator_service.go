package services

import (
	"fmt"

	"github.com/samuelloza/isolate-wrapper/src/application/abstractions"
	"github.com/samuelloza/isolate-wrapper/src/domain"
)

const (
	OJ_WT0 = 0  // Wait 0
	OJ_WT1 = 1  // Wait 1
	OJ_CI  = 2  // Compiling Input Error
	OJ_RI  = 3  // Running Input Error
	OJ_AC  = 4  // Accepted
	OJ_PE  = 5  // Presentation Error
	OJ_WA  = 6  // Wrong Answer
	OJ_TL  = 7  // Time Limit Exceeded
	OJ_ML  = 8  // Memory Limit Exceeded
	OJ_OL  = 9  // Output Limit Exceeded
	OJ_RE  = 10 // Runtime Error
	OJ_CE  = 11 // Compilation Error
	OJ_CO  = 12 // Compiler Output
	OJ_TR  = 13 // Truncated Output
)

type EvaluatorService struct {
	Sandbox    abstractions.Sandbox
	Compiler   abstractions.Compiler
	FileSystem abstractions.FileSystem
	Comparator abstractions.Comparator
}

func NewEvaluatorService(sandbox abstractions.Sandbox, compiler abstractions.Compiler, fileSystem abstractions.FileSystem, comparator abstractions.Comparator) *EvaluatorService {
	return &EvaluatorService{
		Sandbox:    sandbox,
		Compiler:   compiler,
		FileSystem: fileSystem,
		Comparator: comparator,
	}
}

func (s *EvaluatorService) Evaluate(input domain.EvaluationInput) (domain.EvaluationResult, error) {

	evaluationResult := domain.EvaluationResult{
		SubmitID:    input.ID,
		Results:     nil,
		TotalPassed: 0,
		TotalCases:  len(input.TestCases),
		Status:      OJ_CO,
	}

	// Initialize sandbox
	if err := s.Sandbox.Init(input.BoxID); err != nil {
		return evaluationResult, fmt.Errorf("failed to initialize sandbox: %w", err)
	}
	defer s.Sandbox.Cleanup(input.BoxID)
	defer s.FileSystem.DeleteDir(fmt.Sprint(input.BoxID))

	s.FileSystem.CreateTmpDirectory(input.BoxID)
	// Write source code to temporary directory
	srcFileName := fmt.Sprintf("Main.%s", input.Language)
	if err := s.FileSystem.WriteFile(input.BoxID, srcFileName, input.SourceCode); err != nil {
		return evaluationResult, fmt.Errorf("failed to write source code: %w", err)
	}

	// Compile source code
	outputDir := fmt.Sprintf("/tmp/patito-wrapper-%d", input.BoxID)
	if err := s.Compiler.Compile(srcFileName, outputDir); err != nil {
		evaluationResult.Status = OJ_CE
		return evaluationResult, fmt.Errorf("compilation failed: %w", err)
	}

	// Execute test cases
	var results []domain.TestCaseResult
	totalPassed := 0

	for i, test := range input.TestCases {
		// Prepare input and expected output files
		if err := s.FileSystem.CopyFile(test.Input, input.BoxID, "input.txt"); err != nil {
			return evaluationResult, fmt.Errorf("failed to copy input file: %w", err)
		}

		if err := s.FileSystem.CopyFile(test.Output, input.BoxID, "expected.txt"); err != nil {
			return evaluationResult, fmt.Errorf("failed to copy expected output file: %w", err)
		}

		// Run the program in the sandbox
		logData, err := s.Sandbox.Run(input.BoxID, i)
		if err != nil {
			results = append(results, domain.TestCaseResult{
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
			results = append(results, domain.TestCaseResult{
				Index:        i,
				Passed:       false,
				ErrorMessage: err.Error(),
			})
			continue
		}

		// Record result
		passed := cmpResult == domain.OJ_AC
		if passed {
			totalPassed++
		}
		results = append(results, domain.TestCaseResult{
			Index:         i,
			Passed:        passed,
			Expected:      test.Output,
			Received:      outputPath,
			ExecutionTime: logData.ExecutionTime,
			MemoryUsed:    logData.MemoryUsed,
		})
	}

	evaluationResult.Results = results
	evaluationResult.TotalPassed = totalPassed
	evaluationResult.TotalCases = len(results)
	if totalPassed == len(results) {
		evaluationResult.Status = OJ_AC
	} else {
		evaluationResult.Status = OJ_WA
	}

	return evaluationResult, nil
}
