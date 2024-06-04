package main

// import (
// 	"bufio"
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"os"
// 	"path/filepath"

// 	"github.com/jplindgren/bastard-git/object"
// )

// func GenerateObjects(dirPath string) ([]object.BGitObject, error) {
// 	repo := GetRepository()
// 	fmt.Printf("Reading directory: %s\n", dirPath)

// 	entries, err := os.ReadDir(dirPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var objs []object.BGitObject

// 	for _, entry := range entries {
// 		if entry.IsDir() {
// 			if entry.Name() == repo.bGitFolder {
// 				continue
// 			}

// 			generatedObjs, err := GenerateObjects(filepath.Join(dirPath, entry.Name()))
// 			objs = append(objs, generatedObjs...)
// 			if err != nil {
// 				return nil, err
// 			}
// 		} else {
// 			path := filepath.Join(dirPath, entry.Name())
// 			file, err := os.Open(path)

// 			info, _ := file.Stat()

// 			if err != nil {
// 				fmt.Printf("Cannot read file %s\n", file.Name())
// 				os.Exit(1)
// 			}

// 			var out bytes.Buffer
// 			r := bufio.NewReader(file)
// 			io.Copy(&out, r)

// 			fileName := filepath.Base(file.Name())
// 			filePath := file.Name()
// 			bGitBlob := object.NewBlob(out.String(), fileName, info.Mode(), info.ModTime(), info.Size(), filePath)
// 			repo.store.Write(bGitBlob)

// 			objs = append(objs, bGitBlob)

// 			file.Close()
// 		}
// 	}

// 	dirInfo, _ := os.Stat(dirPath)
// 	tree := object.NewTree(objs, dirPath, dirInfo.Mode(), dirInfo.ModTime())
// 	repo.store.Write(tree)
// 	objs = append(objs, tree)

// 	return objs, nil
// }
