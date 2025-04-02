package application

import (
	"testing"

	application "github.com/samuelloza/isolate-wrapper/src/application/services"
	"github.com/samuelloza/isolate-wrapper/src/domain"

	"github.com/davecgh/go-spew/spew"
)

func runTest(t *testing.T, id string, boxID int, code string, expectAllPass bool) {
	evaluator := application.NewEvaluatorService()
	input := domain.EvaluationInput{
		ID:          id,
		UniqID:      id,
		BoxID:       boxID,
		ProblemName: "Suma de Dos Números",
		Language:    "cpp",
		SourceCode:  code,
		MetaPrefix:  id + "-meta",
		RunLimits: domain.RunLimits{
			Time:   2,
			Memory: 65536,
			Output: 1024,
		},
		TestCases: []domain.TestCase{
			{Input: "../../testcases/test_1.in", Output: "../../testcases/test_1.out"},
			{Input: "../../testcases/test_2.in", Output: "../../testcases/test_2.out"},
		},
	}
	result, err := evaluator.Evaluate(input)
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	spew.Dump(result)

	if expectAllPass && result.TotalPassed != len(result.Results) {
		t.Errorf("Expected all test cases to pass, but only %d/%d passed", result.TotalPassed, len(result.Results))
	}
	if !expectAllPass && result.TotalPassed == len(result.Results) {
		t.Errorf("Expected some test cases to fail, but all passed")
	}
}

func Test_Evaluator_AC(t *testing.T) {
	runTest(t, "test-ac", 1, `
		#include <iostream>
		using namespace std;
		int main() {
		    int a;
		    cin>>a;
		    cout<<a*10<<endl;
		    return 0;
		}`, true)
}

func Test_Evaluator_PE(t *testing.T) {
	runTest(t, "test-pe", 2, `
		#include <iostream>
		using namespace std;
		int main() {
		    int a;
		    cin>>a;
		    cout<<a*10;
		    return 0;
		}`, false)
}

func Test_Evaluator_WA(t *testing.T) {
	runTest(t, "test-wa", 3, `
		#include <iostream>
		using namespace std;
		int main() {
		    int a;
		    cin>>a;
		    cout<<a*10+1<<endl;
		    return 0;
		}`, false)
}

func Test_Evaluator_RTE(t *testing.T) {
	runTest(t, "test-rte", 4, `
		#include <iostream>
		using namespace std;
		int main() {
		    int a = 0;
		    cout << 10 / a << endl;
		    return 0;
		}`, false)
}

func Test_Evaluator_TLE(t *testing.T) {
	runTest(t, "test-tle", 5, `
		#include <iostream>
		using namespace std;
		int main() {
		    while (true) {}
		    return 0;
		}`, false)
}

func Test_Evaluator_CE(t *testing.T) {
	runTest(t, "test-ce", 6, `
		#include <iostream>
		int main() {
		    cout << "Hello"
		    return 0;
		}`, false)
}
