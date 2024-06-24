package repository

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jplindgren/bastard-git/internal/object"
	"github.com/jplindgren/bastard-git/internal/utils"
)

const (
	COMMIT = iota
	CHECKOUT
)

type HeadLogType uint8

func (s HeadLogType) String() string {
	switch s {
	case COMMIT:
		return "commit"
	case CHECKOUT:
		return "checkout"
	default:
		return ""
	}
}

func (r *Repository) logHead(logType HeadLogType, parentCommit string, commit string, t time.Time, message string) error {
	file, err := os.OpenFile(r.paths.logs.head, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if parentCommit == "" {
		parentCommit = "0000000000000000000000000000000000000000"
	}

	logEntry := fmt.Sprintf("%s %s <%s> %s %s: %s\n", parentCommit, commit, r.Email, t.Format(time.RFC3339), logType.String(), message)
	_, err = file.WriteString(logEntry)
	return err
}

func (r *Repository) logHeadRef(branchRef string, parentCommit string, commit string, t time.Time, message string) error {
	branchRefPath := filepath.Join(r.paths.logs.path, branchRef)
	if err := utils.CreatePathFoldersIfNotExists(branchRefPath); err != nil {
		return err
	}

	file, err := os.OpenFile(branchRefPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if parentCommit == "" {
		parentCommit = "0000000000000000000000000000000000000000"
	}

	logEntry := fmt.Sprintf("%s %s <%s> %s %s: %s\n", parentCommit, commit, r.Email, t.Format(time.RFC3339), "commit", message)
	_, err = file.WriteString(logEntry)
	return err
}

func (r *Repository) GetLogsForCurrentBranch() error {
	commitHash, _, err := r.lookupLastCommit()
	if err != nil {
		return err
	}

	for commitHash != "" {
		cContent, err := r.Store.Get(commitHash)
		if err != nil {
			return err
		}

		cContent, found := bytes.CutPrefix(cContent, object.CommitPrefix)
		if !found {
			return fmt.Errorf("could not parse object as commit: %s", commitHash)
		}

		sContent := string(cContent)

		fmt.Fprintf(os.Stdout, "commit %s\n%s\n\n*********************\n", commitHash, sContent)

		lines := strings.Split(sContent, "\n")
		parentLine := strings.Split(lines[1], " ")
		commitHash = parentLine[1]
	}

	return nil
}
