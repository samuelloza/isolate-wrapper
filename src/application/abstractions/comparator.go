package abstractions

type Comparator interface {
	Compare(expectedPath string, outputPath string) (ComparisonResult, error)
}

type ComparisonResult int

const (
	OJ_WT0 = 0  // Wait 0
	OJ_WT1 = 1  // Wait 1
	OJ_CI  = 2  // Compiling Input Error
	OJ_RI  = 3  // Running Input Error
	OJ_AC  = 4  // Accepted
	OJ_PE  = 5  // Presentation Error
	OJ_WA  = 6  // Wrong Answer
	OJ_TL  = 7  // Time Limit Exceeded
	OJ_ML  = 8  // Memory Limit Exceeded
	OJ_OL  = 9  // Output Limit Exceeded
	OJ_RE  = 10 // Runtime Error
	OJ_CE  = 11 // Compilation Error
	OJ_CO  = 12 // Compiler Output
	OJ_TR  = 13 // Truncated Output
)
