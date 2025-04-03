package compiler

import (
	"os/exec"
)

type PseintCompiler struct {
	SandBoxDir string
}

func (p *PseintCompiler) Compile(_ string, boxDir string) error {
	shCmd := "chown judge:judge Main.psc && /usr/bin/dos2unix -b Main.psc && /usr/bin/pseint Main.psc --draw Main.psd --fixwincharset --norun pseint.txt && /usr/bin/psexport --lang=cpp Main.psd Main.cc && g++ Main.cc -o Main -fno-asm -Wall -lm --static -DONLINE_JUDGE"
	cmd := exec.Command("/bin/sh", "-c", shCmd)
	cmd.Dir = boxDir
	return cmd.Run()
}
