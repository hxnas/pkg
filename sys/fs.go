package sys

import (
	"io"
	"os"
)

func IsFile(path string) bool {
	f, err := os.Stat(path)
	return err == nil && f.Mode().IsRegular()
}

func IsDir(path string) bool {
	f, err := os.Stat(path)
	return err == nil && f.IsDir()
}

func IsDirEmpty(path string) bool {
	_, err := ReadDirNames(path, 1)
	return err == io.EOF
}

func FileNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}
