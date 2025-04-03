package compiler

import (
	"fmt"
)

type Compiler interface {
	Compile(srcPath string, boxDir string) error
}

func GetCompiler(language string, boxID int) (Compiler, error) {
	sandBoxTmpDir := fmt.Sprintf("/tmp/patito-wrapper-%d", boxID)
	switch language {
	case "cpp":
		return &CppCompiler{SandBoxDir: sandBoxTmpDir}, nil
	case "python":
		return &PythonCompiler{SandBoxDir: sandBoxTmpDir}, nil
	case "java":
		return &JavaCompiler{SandBoxDir: sandBoxTmpDir}, nil
	case "pseint":
		return &PseintCompiler{SandBoxDir: sandBoxTmpDir}, nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
}
