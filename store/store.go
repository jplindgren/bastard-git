package store

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jplindgren/bastard-git/object"
)

type Store interface {
	Write(obj object.BGitObject) error
	Get(hash string) ([]byte, error)
}

type FileSystemStore struct {
	rootPath string
}

func New(path string) *FileSystemStore {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Cannot get absolute path of %s\n", path)
		os.Exit(1)
	}

	return &FileSystemStore{rootPath: absPath}
}

func (f *FileSystemStore) Write(obj object.BGitObject) error {

	dirPath, filePath := f.GetPath(obj.GetHash())
	fmt.Printf("Writing object %s of type: %s at %s\n", object.GetSha1AsString(obj.GetHash()), obj.GetType(), filePath)

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	w := zlib.NewWriter(file)
	content := obj.Serialize()
	_, err = w.Write(content)
	w.Close()
	file.Close()

	return err
}

func (f *FileSystemStore) Get(hash string) ([]byte, error) {
	dirName := hash[:2]
	fileName := hash[2:]

	path := filepath.Join(f.rootPath, dirName, fileName)

	objFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer objFile.Close()

	var out bytes.Buffer
	r, err := zlib.NewReader(objFile)
	if err != nil {
		return nil, err
	}

	io.Copy(&out, r)
	return out.Bytes(), nil
}

func (f *FileSystemStore) GetPath(hash []byte) (dir string, file string) {
	prefixName := hex.EncodeToString(hash[:1])
	objName := hex.EncodeToString(hash[1:])

	return filepath.Join(f.rootPath, prefixName), filepath.Join(f.rootPath, prefixName, objName)
}
