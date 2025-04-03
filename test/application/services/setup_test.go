package services_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/samuelloza/isolate-wrapper/src/application/services"
	"github.com/samuelloza/isolate-wrapper/src/domain"
)

var boxPool *services.BoxPool

func TestMain(m *testing.M) {
	boxPool = services.NewBoxPool(2)
	code := m.Run()
	os.Exit(code)
}

func getTestCases(paths []struct{ in, out string }, cwd string) []domain.TestCase {
	var tcs []domain.TestCase
	for _, p := range paths {
		tcs = append(tcs, domain.TestCase{
			Input:  filepath.Join(cwd, p.in),
			Output: filepath.Join(cwd, p.out),
		})
	}
	return tcs
}
