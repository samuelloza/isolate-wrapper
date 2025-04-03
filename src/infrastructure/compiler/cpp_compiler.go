package compiler

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

type CppCompiler struct {
	SandBoxDir string
}

func (c *CppCompiler) Compile(srcPath string, _ string) error {
	outPath := filepath.Join(c.SandBoxDir, "Main")
	cmd := exec.Command("g++", fmt.Sprintf("%s/%s", c.SandBoxDir, srcPath), "-o", outPath, "-fno-asm", "-Wall", "-lm", "--static", "-std=c++11", "-DONLINE_JUDGE")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("compilation failed %s:\n%s\n", srcPath, string(output))
		return err
	}

	return nil
}
