package repository

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ADD = iota
	REMOVE
	MODIFY
)

type Status uint8

func (s Status) String() string {
	switch s {
	case ADD:
		return "add"
	case REMOVE:
		return "remove"
	case MODIFY:
		return "modify"
	default:
		return ""
	}
}

type DiffResult struct {
	Name   string
	Status Status
}

func (r *Repository) Diff() ([]DiffResult, error) {
	toBeCommited := []DiffResult{}

	ignoredPaths, _ := r.getIgnoredPaths()
	indexFiles, err := r.readIndex()
	if err != nil {
		return nil, err
	}

	filepath.Walk(r.WorkTree, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || strings.Contains(path, r.BGitFolder) || strings.Contains(path, bGitIgnoreFile) {
			return nil
		}

		if _, ok := ignoredPaths[path]; ok {
			return nil
		}

		relPath, err := filepath.Rel(r.WorkTree, path)
		if err != nil {
			return err
		}

		found := false
		for _, idxFile := range indexFiles {
			if idxFile.name == relPath {
				found = true
				if info.ModTime().Format(time.RFC3339) != idxFile.modTime {
					toBeCommited = append(toBeCommited, DiffResult{Name: info.Name(), Status: MODIFY})
				}
			}
		}

		if !found {
			toBeCommited = append(toBeCommited, DiffResult{Name: info.Name(), Status: ADD})
		}

		return nil
	})

	return toBeCommited, nil
}
