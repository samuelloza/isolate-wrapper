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

const (
	OJ_AC  = 0
	OJ_WA  = 1
	OJ_TLE = 2
	OJ_MLE = 3
	OJ_RE  = 4
	OJ_CE  = 5
	OJ_PE  = 6
)
