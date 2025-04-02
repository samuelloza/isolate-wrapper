package abstractions

type Comparator interface {
	Compare(expectedPath string, outputPath string) (ComparisonResult, error)
}

type ComparisonResult int

const (
	OJ_AC  ComparisonResult = iota // Accepted
	OJ_WA                          // Wrong Answer
	OJ_PE                          // Presentation Error
	OJ_TLE                         // Time Limit Exceeded
	OJ_MLE                         // Memory Limit Exceeded
	OJ_RE                          // Runtime Error
	OJ_CE                          // Compilation Error
)
