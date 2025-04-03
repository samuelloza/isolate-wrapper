package compiler

import (
	"fmt"
	"os/exec"
)

type PythonCompiler struct {
	SandBoxDir string
}

func (p *PythonCompiler) Compile(srcPath string, _ string) error {
	cmd := exec.Command("/usr/bin/python3.12", "-c", fmt.Sprintf("import py_compile; py_compile.compile(r'%s')", srcPath))
	return cmd.Run()
}
