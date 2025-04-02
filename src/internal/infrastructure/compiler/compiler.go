package compiler

import (
	"fmt"
)

type Compiler interface {
	Compile(srcPath string, boxDir string) error
}

func GetCompiler(language string, SandBoxDir string) (Compiler, error) {
	switch language {
	case "cpp":
		return &CppCompiler{SandBoxDir: SandBoxDir}, nil
	case "python":
		return &PythonCompiler{SandBoxDir: SandBoxDir}, nil
	case "java":
		return &JavaCompiler{SandBoxDir: SandBoxDir}, nil
	case "pseint":
		return &PseintCompiler{SandBoxDir: SandBoxDir}, nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
}
