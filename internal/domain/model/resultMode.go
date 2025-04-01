package model

type TestCaseResult struct {
	Index         int
	Passed        bool
	Expected      string
	Received      string
	ExecutionTime int
	MemoryUsed    int
	ErrorMessage  string
}

type EvaluationResult struct {
	EvaluationID string
	Results      []TestCaseResult
	TotalPassed  int
	TotalCases   int
	Status       string
}
