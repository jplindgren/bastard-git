package repository

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jplindgren/bastard-git/internal/utils"
)

func (r *Repository) DeleteWorkingTree() error {
	dirEntries, err := os.ReadDir(r.WorkTree)
	if err != nil {
		return err
	}

	for _, entry := range dirEntries {
		if entry.Name() == r.BGitFolder || entry.Name() == r.bGitTempFolder {
			continue
		}

		err := os.RemoveAll(filepath.Join(r.WorkTree, entry.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

func RecreateWorkingTree(treeHash string) error {
	var idxBuffer bytes.Buffer
	repo := GetRepository("")

	// Create temp folder
	err := os.Mkdir(repo.Paths.bGitTempPath, 0755)
	if err != nil {
		return err
	}

	// Use temp folder to rollback in case it fails
	err = repo.RecreateWorkingTree(treeHash, repo.Paths.bGitTempPath, &idxBuffer)
	if err != nil {
		return err
	}

	err = repo.DeleteWorkingTree()
	if err != nil {
		// Remove temp folder with new working tree if deleting the current one fails
		os.RemoveAll(repo.Paths.bGitTempPath)
		return err
	}

	//After use temp folder, we should copy content from it to working tree and delete it
	err = utils.CopyDir(repo.WorkTree, repo.Paths.bGitTempPath)
	if err != nil {
		return err
	}

	// Remove temp file
	err = os.RemoveAll(repo.Paths.bGitTempPath)
	if err != nil {
		return err
	}

	// Write updated index
	return repo.updateIndex(&idxBuffer)
}

func (r *Repository) RecreateWorkingTree(treeHash string, path string, index *bytes.Buffer) error {
	content, err := r.Store.Get(treeHash)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 4 {
			return fmt.Errorf("invalid tree entry: %s", line)
		}

		if parts[1] == "tree" {
			err = os.Mkdir(filepath.Join(path, parts[3]), 0755)
			if err != nil {
				return err
			}

			err = r.RecreateWorkingTree(parts[2], filepath.Join(path, parts[3]), index)
			if err != nil {
				return err
			}
		} else {
			fileHash := parts[2]
			fileContent, err := r.Store.Get(fileHash)
			if err != nil {
				return err
			}

			tempFilePath := filepath.Join(path, parts[3])
			newFile, err := os.Create(tempFilePath) //TODO: create file with correct permissions based on tree object
			if err != nil {
				return err
			}
			defer newFile.Close()

			_, err = newFile.Write(fileContent)
			if err != nil {
				return err
			}

			//TODO: should not use bGitTempFolder directly here.
			indexFilePath := strings.ReplaceAll(tempFilePath, fmt.Sprintf("%s/", bGitTempFolder), "")
			info, _ := newFile.Stat()
			index.WriteString(fmt.Sprintf("%s %d %s %s\n", fileHash, len(fileContent), indexFilePath, info.ModTime().Format(time.RFC3339)))
		}
	}

	return nil
}
