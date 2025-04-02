package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/samuelloza/isolate-wrapper/src/internal/application"
	"github.com/samuelloza/isolate-wrapper/src/internal/domain/model"
)

func main() {
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
			{Input: "../../example/test_1.in", Output: "../../example/test_1.out"},
			{Input: "../../example/test_2.in", Output: "../../example/test_2.out"},
		},
	}

	evaluator := application.NewEvaluatorService()
	result, err := evaluator.Evaluate(input)
	if err != nil {
		log.Fatalf("Error evaluating: %v", err)
	}

	spew.Dump(result)
}
