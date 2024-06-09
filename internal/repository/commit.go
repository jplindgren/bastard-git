package repository

import (
	"bytes"
	"path/filepath"

	"github.com/jplindgren/bastard-git/internal/object"
	"github.com/jplindgren/bastard-git/internal/utils"
)

func (r *Repository) Commit(message string) error {
	var buffer bytes.Buffer
	rootTree, err := r.writeTree(r.WorkTree, &buffer)
	if err != nil {
		return err
	}

	lastCommit, err := r.lookupLastCommit()
	if err != nil {
		return err
	}

	commit := object.NewCommit(r.Email, message, rootTree.GetHash(), []string{lastCommit})
	err = r.Store.Write(commit)
	if err != nil {
		return err
	}

	r.updateIndex(&buffer)
	r.UpdateRefHead(commit.ToString())
	return nil
}

func (r *Repository) lookupLastCommit() (string, error) {
	//TODO: commit should check if a previous commit exists and set it as parent
	branch, err := r.GetCurrentBranch()
	if err != nil {
		return "", err
	}

	commitHash, _, err := r.GetBranchTip(branch)
	if err != nil {
		return "", err
	}

	return commitHash, nil
}

// TODO: change indexEntries to be a struct? or keep as string?
// writeTree still loads all file contents to memory. Maybe change store.Write to receive the content via parameter, and keep the object trees lean?
func (r *Repository) writeTree(dirPath string, indexEntries *bytes.Buffer) (*object.Tree, error) {
	entries, dirInfo, err := utils.ReadDir(dirPath, true)
	if err != nil {
		return nil, err
	}

	var objs []object.BGitObject
	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == r.BGitFolder { //ignore special folder
				continue
			}

			innerTree, err := r.writeTree(filepath.Join(dirPath, entry.Name()), indexEntries)
			if err != nil {
				return nil, err
			}

			objs = append(objs, innerTree)
		} else {
			path := filepath.Join(dirPath, entry.Name())

			fBuffer, fInfo, err := utils.ReadFile(path, true)
			if err != nil {
				return nil, err
			}

			fileName := filepath.Base(fInfo.Name())
			filePath := fInfo.Name()
			bGitBlob := object.NewBlob(fBuffer.String(), fileName, fInfo.Mode(), fInfo.ModTime(), fInfo.Size(), filePath)
			r.Store.Write(bGitBlob)

			objs = append(objs, bGitBlob)
		}
	}

	tree := object.NewTree(objs, dirPath, dirInfo.Mode(), dirInfo.ModTime())
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
