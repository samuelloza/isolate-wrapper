package abstractions

type FileSystem interface {
	WriteFile(boxID int, fileName string, content string) error
	CopyFile(srcPath string, boxID int, destFileName string) error
	GetFilePath(boxID int, fileName string) string
	GetOutputPath(boxID int) string
	GetErrorPath(boxID int) string
	DeleteDir(s string) error
	CreateTmpDirectory(boxID int) error
}
