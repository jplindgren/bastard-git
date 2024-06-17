package repository

import (
	"bytes"
	"path/filepath"

	"github.com/jplindgren/bastard-git/internal/object"
	"github.com/jplindgren/bastard-git/internal/utils"
)

func (r *Repository) Commit(message string) error {
	var buffer bytes.Buffer

	ignoredPaths, _ := r.getIgnoredPaths()
	rootTree, err := r.writeTree(r.WorkTree, &buffer, ignoredPaths)
	if err != nil {
		return err
	}

	lastCommit, branchRef, err := r.lookupLastCommit()
	if err != nil {
		return err
	}

	commit := object.NewCommit(r.Email, message, rootTree.GetHash(), []string{lastCommit})
	err = r.Store.Write(commit)
	if err != nil {
		return err
	}

	err = r.updateIndex(&buffer)
	if err != nil {
		return err
	}

	//TODO: rework logs. updateRefHead knows the refHead, and should be able to log the commit.
	//Also, both logs may have the same contract
	err = r.updateRefHead(commit.ToString(), branchRef)
	if err != nil {
		return err
	}

	err = r.logHeadRef(branchRef, lastCommit, commit.ToString(), commit.CreatedAt, message)
	if err != nil {
		return err
	}

	err = r.logHead(COMMIT, lastCommit, commit.ToString(), commit.CreatedAt, message)
	return err
}

// Returns the last commit hash and branch ref
func (r *Repository) lookupLastCommit() (string, string, error) {
	_, branchRef, err := r.GetCurrentBranch()
	if err != nil {
		return "", "", err
	}

	commitHash, _, err := r.GetBranchTip(branchRef)
	if err != nil {
		return "", "", err
	}

	return commitHash, branchRef, nil
}

// TODO: change indexEntries to be a struct? or keep as string?
// writeTree still loads all file contents to memory. Maybe change store.Write to receive the content via parameter, and keep the object trees lean?
func (r *Repository) writeTree(treePath string, indexEntries *bytes.Buffer, ignoredPaths map[string]bool) (*object.Tree, error) {
	entries, dirInfo, err := utils.ReadDir(treePath, true)
	if err != nil {
		return nil, err
	}

	var objs []object.BGitObject
	for _, entry := range entries {

		absEntryPath := filepath.Join(treePath, entry.Name())
		if _, ok := ignoredPaths[absEntryPath]; ok {
			continue
		}

		if entry.IsDir() {
			if entry.Name() == r.BGitFolder { //ignore bgit folder
				continue
			}

			innerTree, err := r.writeTree(absEntryPath, indexEntries, ignoredPaths)
			if err != nil {
				return nil, err
			}

			if innerTree == nil {
				continue
			}

			objs = append(objs, innerTree)
		} else {
			fBuffer, fInfo, err := utils.ReadFile(absEntryPath, true)
			if err != nil {
				return nil, err
			}

			relFilePath, err := filepath.Rel(r.WorkTree, absEntryPath)
			if err != nil {
				return nil, err
			}

			bGitBlob := object.NewBlob(fBuffer.String(), entry.Name(), fInfo.Mode(), fInfo.ModTime(), fInfo.Size(), relFilePath)
			r.Store.Write(bGitBlob)

			objs = append(objs, bGitBlob)
		}
	}

	if len(objs) == 0 { // ignore trees with no objects
		return nil, nil
	}

	tree := object.NewTree(objs, treePath, dirInfo.Mode(), dirInfo.ModTime())
	err = r.Store.Write(tree)
	if err != nil {
		return nil, err
	}

	_, err = indexEntries.WriteString(tree.FormatToIndex())
	if err != nil {
		return nil, err
	}

	return tree, nil
}
