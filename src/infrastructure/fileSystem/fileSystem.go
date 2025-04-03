package fileSystem

import (
	"fmt"
	"io"
	"os"
)

type FileSystem struct{}

func (fs *FileSystem) WriteFile(boxID int, fileName string, content string) error {
	dst := fmt.Sprintf("/tmp/patito-wrapper-%d/%s", boxID, fileName)
	return os.WriteFile(dst, []byte(content), 0644)
}

func (fs *FileSystem) CopyFile(srcPath string, boxID int, destFileName string) error {
	dst := fmt.Sprintf("/tmp/patito-wrapper-%d/%s", boxID, destFileName)
	return copyFile(srcPath, dst)
}

func (fs *FileSystem) GetFilePath(boxID int, fileName string) string {
	return fmt.Sprintf("/tmp/patito-wrapper-%d/%s", boxID, fileName)
}

func (fs *FileSystem) GetOutputPath(boxID int) string {
	return fmt.Sprintf("/var/local/lib/isolate/%d/box/user_output.txt", boxID)
}

func (fs *FileSystem) GetErrorPath(boxID int) string {
	return fmt.Sprintf("/var/local/lib/isolate/%d/box/error.txt", boxID)
}

func (fs *FileSystem) CreateTmpDirectory(boxID int) error {
	dst := fmt.Sprintf("/tmp/patito-wrapper-%d", boxID)
	return CreateDirectory(dst)
}

func (fs *FileSystem) DeleteDir(path string) error {
	return os.RemoveAll(fmt.Sprintf("/tmp/patito-wrapper-%s", path))
}

func CreateDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
