package abstractions

import "github.com/samuelloza/isolate-wrapper/src/domain"

type TestCaseProvider interface {
	GetTestCases(problemID string) ([]domain.TestCase, error)
}
