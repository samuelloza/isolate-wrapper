package abstractions

type Sandbox interface {
	Init(boxID int) error
	Run(boxID int, testCaseIndex int) (SandboxLogData, int)
	Cleanup(boxID int) error
}

type SandboxLogData struct {
	ExecutionTime float64
	MemoryUsed    int
}
