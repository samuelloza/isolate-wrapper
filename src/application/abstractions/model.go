package abstractions

type IsolateSandbox struct {
	BoxID            int
	RunLogPrefix     string
	Chdir            string
	PreserveEnv      bool
	InheritEnv       []string
	SetEnv           map[string]string
	FSize            int
	StdinFile        string
	StackSpace       int
	AddressSpace     int
	StdoutFile       string
	StderrFile       string
	Timeout          float64
	WallclockTimeout float64
	ExtraTimeout     float64
	Verbosity        int
}
