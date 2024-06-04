package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jplindgren/bastard-git/object"
)

type FileEntry struct {
	size    int
	name    string
	path    string
	isDir   bool
	mode    os.FileMode
	content string
}

type FileReader interface {
	Read() ([]FileEntry, error)
}

type FileSystemReader struct {
}

func (fst *FileSystemReader) GetAbsolutePath(path string) string {
	absDir, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Cannot get absolute path of %s\n", path)
		os.Exit(1)

	}
	return absDir
}

func (fsr *FileSystemReader) Read(dirPath string) ([]FileEntry, error) {
	fmt.Printf("Reading directory: %s\n", dirPath)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var fileEntries []FileEntry

	for _, entry := range entries {
		if entry.IsDir() {
			fileEntries = append(fileEntries, FileEntry{
				size:  0,
				name:  entry.Name(),
				isDir: true,
				path:  filepath.Join(dirPath, entry.Name()),
				mode:  entry.Type()})
		} else {
			path := filepath.Join(dirPath, entry.Name())
			file, err := os.Open(path)

			if err != nil {
				fmt.Printf("Cannot read file %s\n", file.Name())
				os.Exit(1)
			}

			var out bytes.Buffer
			r := bufio.NewReader(file)
			io.Copy(&out, r)

			data := out.String()

			info, _ := file.Stat()
			fileEntries = append(fileEntries, FileEntry{
				size:    int(info.Size()),
				name:    entry.Name(),
				path:    path,
				isDir:   false,
				mode:    info.Mode(),
				content: data,
			})
			file.Close()
		}
	}

	return fileEntries, nil
}

func (fsr *FileSystemReader) WalkFiles(dirPath string) {
	counter := 0
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		counter = counter + 1

		fmt.Printf("%s %s %o\n", info.Name(), path, info.Mode())
		return nil
	})

	fmt.Println("Number of items -> " + fmt.Sprint(counter))
}

func (fsr *FileSystemReader) Walk(dirPath string) {
	repo := GetRepository()
	counter := 0

	filepath.WalkDir(dirPath, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, repo.bGitFolder) {
			return nil
		}

		info, _ := dir.Info()
		fmt.Printf("%s %s %d \n", dir.Name(), path, info.Mode())

		counter++

		if !dir.IsDir() {
			entries, err := os.ReadDir(path)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				// fmt.Printf("Reading file: %s %s\n", entry.Name(), filepath.Join(path, entry.Name()))
				// fileBytes, err := os.ReadFile(filepath.Join(path, entry.Name()))
				// if err != nil {
				// 	fmt.Println(err)
				// 	os.Exit(1)
				// }

				// bGitObject := object.NewBlob(string(fileBytes), len(fileBytes))
				// repo.store.Write(bGitObject)
			}

		} else {
			// create tree
			return nil
		}

		return nil
	})
	fmt.Println("Number of items -> " + fmt.Sprint(counter))
}

type DiffResult struct {
	name   string
	status string //add, remove, modify
}

func (fsr *FileSystemReader) Diff(rootPath string) ([]DiffResult, error) {
	repo := GetRepository()

	toBeCommited := []DiffResult{}

	indexFiles, err := readIndex(repo.paths.indexPath)
	if err != nil {
		return nil, err
	}

	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || strings.Contains(path, repo.bGitFolder) {
			return nil
		}

		found := false
		for _, idxFile := range indexFiles {
			if idxFile.name == path {
				found = true
				if info.ModTime().Format(time.RFC3339) != idxFile.modTime {
					toBeCommited = append(toBeCommited, DiffResult{name: info.Name(), status: "modify"})
				}
			}
		}

		if !found {
			toBeCommited = append(toBeCommited, DiffResult{name: info.Name(), status: "add"})
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
	repo := GetRepository()
	fmt.Printf("Reading directory: %s\n", dirPath)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var objs []object.BGitObject

	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == repo.bGitFolder {
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
			repo.store.Write(bGitBlob)

			objs = append(objs, bGitBlob)

			file.Close()
		}
	}

	dirInfo, _ := os.Stat(dirPath)
	tree := object.NewTree(objs, dirPath, dirInfo.Mode(), dirInfo.ModTime())
	repo.store.Write(tree)

	return tree, nil
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
