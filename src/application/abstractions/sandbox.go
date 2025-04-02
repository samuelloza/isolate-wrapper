package abstractions

type Sandbox interface {
	Init(boxID int) error
	Run(boxID int, testCaseIndex int) (SandboxLogData, error)
	Cleanup(boxID int) error
}

type SandboxLogData struct {
	ExecutionTime int
	MemoryUsed    int
}
