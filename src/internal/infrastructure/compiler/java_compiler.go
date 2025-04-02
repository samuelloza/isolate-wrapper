package compiler

import (
	"os/exec"
)

type JavaCompiler struct {
	SandBoxDir string
}

func (j *JavaCompiler) Compile(srcPath string, boxDir string) error {
	cmd := exec.Command("javac", "-J-Xms32m", "-J-Xmx256m", "-encoding", "UTF-8", srcPath)
	cmd.Dir = boxDir
	return cmd.Run()
}
