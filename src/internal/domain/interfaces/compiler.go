package interfaces

type Compiler interface {
	Compile(srcPath string, boxDir string) error
}
