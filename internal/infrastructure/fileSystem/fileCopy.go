package fileSystem

import (
	"fmt"
	"io"
	"os"
)

func WriteCodeInTmpDirectory(code string, boxID int, fileName string) error {
	dst := fmt.Sprintf("/tmp/patito-wrapper-%d/%s", boxID, fileName)
	err := os.WriteFile(dst, []byte(code), 0644)
	if err != nil {
		return err
	}
	return nil
}

func CopyFileToTmpDirectory(srcPath string, boxID int, fileName string) error {
	dst := fmt.Sprintf("/tmp/patito-wrapper-%d/%s", boxID, fileName)
	return copyFile(srcPath, dst)
}

func CreateTmpDirectory(boxID int) {
	dst := fmt.Sprintf("/tmp/patito-wrapper-%d", boxID)
	CreateDirectory(dst)
}

func DeleteDirectoryFromTmpDirectory(boxID int, fileName string) {
	dst := fmt.Sprintf("/tmp/patito-wrapper-%d/%s", boxID, fileName)
	CreateDirectory(dst)
}

func CreateDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func DeleteDir(path string) error {
	return os.RemoveAll(path)
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
