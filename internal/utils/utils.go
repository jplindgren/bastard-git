package utils

import (
	"bufio"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ReadDir(path string, readInfo bool) ([]fs.DirEntry, os.FileInfo, error) {
	var info os.FileInfo

	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}

	if readInfo {
		info, err = os.Stat(path)
		if err != nil {
			return nil, nil, err
		}
	}

	return dirEntries, info, nil
}

func ReadFile(path string, readInfo bool) (*bytes.Buffer, os.FileInfo, error) {
	var info os.FileInfo

	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	if readInfo {
		info, err = file.Stat()
		if err != nil {
			return nil, nil, err
		}
	}

	var out bytes.Buffer
	r := bufio.NewReader(file)
	io.Copy(&out, r)

	return &out, info, nil
}

// CopyDir copies the content of src to dst. src should be a full path.
func CopyDir(dst, src string) error {

	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// copy to this path
		outpath := filepath.Join(dst, strings.TrimPrefix(path, src))

		if info.IsDir() {
			os.MkdirAll(outpath, info.Mode())
			return nil // means recursive
		}

		// copy contents of regular file efficiently

		// open input
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		// create output
		fh, err := os.Create(outpath)
		if err != nil {
			return err
		}
		defer fh.Close()

		// make it the same
		fh.Chmod(info.Mode())

		// copy content
		_, err = io.Copy(fh, in)
		return err
	})
}

func CheckIfFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CreatePathFoldersIfNotExists(path string) error {
	if !CheckIfFileExists(path) {
		dirPath := filepath.Dir(path)
		err := os.MkdirAll(dirPath, fs.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
