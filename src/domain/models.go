package domain

type RunLimits struct {
	Time   int `json:"time"`
	Memory int `json:"memory"`
	Output int `json:"output"`
}

type TestCase struct {
	Type   string `json:"type"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

type EvaluationInput struct {
	ID          string     `json:"id"`
	UniqID      string     `json:"uniq_id"`
	ProblemName string     `json:"problem_name"`
	Language    string     `json:"language"`
	SourceCode  string     `json:"source_code_path"`
	BoxID       int        `json:"box_id"`
	RunLimits   RunLimits  `json:"run_limits"`
	TestCases   []TestCase `json:"testcases"`
}

type SandboxLog struct {
	ExecutionTime int
	MemoryUsed    int
}
