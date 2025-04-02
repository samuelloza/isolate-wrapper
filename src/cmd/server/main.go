package main

import (
	"log"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/samuelloza/isolate-wrapper/src/internal/application"
	"github.com/samuelloza/isolate-wrapper/src/internal/domain/model"
	"github.com/samuelloza/isolate-wrapper/src/internal/infrastructure/comparator"
	"github.com/samuelloza/isolate-wrapper/src/internal/infrastructure/compiler"
	"github.com/samuelloza/isolate-wrapper/src/internal/infrastructure/fileSystem"
	"github.com/samuelloza/isolate-wrapper/src/internal/infrastructure/isolate"
)

func main() {
	cwd := "/home/sam/project/github/isolate-wrapper"

	input := model.EvaluationInput{
		ID:          "123456",
		UniqID:      "abc123",
		BoxID:       1,
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
		MetaPrefix: "abc123-meta",
		RunLimits: model.RunLimits{
			Time:   2,
			Memory: 65536,
			Output: 1024,
		},
		TestCases: []model.TestCase{
			{Input: filepath.Join(cwd, "test/example/test_1.in"), Output: filepath.Join(cwd, "test/example/test_1.out")},
			{Input: filepath.Join(cwd, "test/example/test_2.in"), Output: filepath.Join(cwd, "test/example/test_2.out")},
		},
	}

	sandboxImpl := &isolate.IsolateSandbox{}
	compilerImpl, err := compiler.GetCompiler(input.Language, "/tmp/patito-wrapper-1")
	if err != nil {
		log.Fatalf("Compiler error: %v", err)
	}
	fileSystemImpl := &fileSystem.FileSystem{}
	comparatorImpl := &comparator.Comparator{}

	evaluator := application.NewEvaluatorService(sandboxImpl, compilerImpl, fileSystemImpl, comparatorImpl)

	result, err := evaluator.Evaluate(input)
	if err != nil {
		log.Fatalf("Error evaluating: %v", err)
	}

	spew.Dump(result)
}
