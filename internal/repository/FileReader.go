package repository

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jplindgren/bastard-git/internal/object"
)

type FileSystemReader struct {
}

type DiffResult struct {
	Name   string
	Status string //add, remove, modify
}

func (fsr *FileSystemReader) Diff(rootPath string) ([]DiffResult, error) {
	repo := GetRepository("")

	toBeCommited := []DiffResult{}

	indexFiles, err := readIndex(repo.Paths.IndexPath)
	if err != nil {
		return nil, err
	}

	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || strings.Contains(path, repo.BGitFolder) {
			return nil
		}

		found := false
		for _, idxFile := range indexFiles {
			if idxFile.name == path {
				found = true
				if info.ModTime().Format(time.RFC3339) != idxFile.modTime {
					toBeCommited = append(toBeCommited, DiffResult{Name: info.Name(), Status: "modify"})
				}
			}
		}

		if !found {
			toBeCommited = append(toBeCommited, DiffResult{Name: info.Name(), Status: "add"})
		}

		return nil
	})

	return toBeCommited, nil
}

type IndexEntry struct {
	name    string
	size    int64
	modTime string
}

func readIndex(indexPath string) ([]IndexEntry, error) {
	indexEntries := []IndexEntry{}

	file, err := os.Open(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return indexEntries, nil
		}
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fLine := scanner.Text()
		data := strings.Split(fLine, " ")
		size, err := strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			return nil, err
		}

		modTime, err := time.Parse(time.RFC3339, data[3])
		if err != nil {
			return nil, err
		}

		indexEntries = append(indexEntries, IndexEntry{
			size:    size,
			name:    data[2],
			modTime: modTime.Format(time.RFC3339),
		})
	}

	return indexEntries, nil
}

func GenerateObjectTree(dirPath string) (object.BGitObject, error) {
	repo := GetRepository("")
	fmt.Printf("Reading directory: %s\n", dirPath)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var objs []object.BGitObject

	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == repo.BGitFolder {
				continue
			}

			innerTree, err := GenerateObjectTree(filepath.Join(dirPath, entry.Name()))
			objs = append(objs, innerTree)
			if err != nil {
				return nil, err
			}
		} else {
			path := filepath.Join(dirPath, entry.Name())
			file, err := os.Open(path)

			info, _ := file.Stat()

			if err != nil {
				fmt.Printf("Cannot read file %s\n", file.Name())
				os.Exit(1)
			}

			var out bytes.Buffer
			r := bufio.NewReader(file)
			io.Copy(&out, r)

			fileName := filepath.Base(file.Name())
			filePath := file.Name()
			bGitBlob := object.NewBlob(out.String(), fileName, info.Mode(), info.ModTime(), info.Size(), filePath)
			repo.Store.Write(bGitBlob)

			objs = append(objs, bGitBlob)

			file.Close()
		}
	}

	dirInfo, _ := os.Stat(dirPath)
	tree := object.NewTree(objs, dirPath, dirInfo.Mode(), dirInfo.ModTime())
	repo.Store.Write(tree)

	return tree, nil
}
