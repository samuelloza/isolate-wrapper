package domain

type EvaluationInput struct {
	ID         string     `json:"id"`
	UniqID     string     `json:"uniqId"`
	BoxID      int        `json:"boxId"`
	ProblemID  int        `json:"problemID"`
	Language   string     `json:"language"`
	SourceCode string     `json:"sourceCode"`
	RunLimits  RunLimits  `json:"runLimits"`
	TestCases  []TestCase `json:"testCases"`
}

type RunLimits struct {
	Time   int `json:"time"`
	Memory int `json:"memory"`
	Output int `json:"output"`
}

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type SandboxLog struct {
	ExecutionTime int
	MemoryUsed    int
}
