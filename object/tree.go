package object

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"time"
)

type Tree struct {
	name    string
	hash    []byte
	entries []BGitObject
	mode    fs.FileMode
	modTime time.Time
	path    string
}

func NewTree(entries []BGitObject, path string, mode fs.FileMode, modTime time.Time) *Tree {
	var data []byte
	for _, entry := range entries {
		data = append(data, entry.GetHash()...)
	}
	hash := generateSHA1Hash(string(data))

	tree := &Tree{
		hash:    hash,
		entries: entries,
		mode:    mode,
		name:    filepath.Base(path),
		path:    path,
		modTime: modTime,
	}

	fmt.Printf("Creating new tree with id %s with entries \n", GetSha1AsString(tree.hash))
	return tree
}

func (t *Tree) Serialize() []byte {
	var buffer bytes.Buffer
	//buffer.WriteString("t\000")
	for _, entry := range t.entries {
		buffer.WriteString(entry.ToString())
	}
	return buffer.Bytes()
}

func (t *Tree) GetHash() []byte {
	return t.hash
}

func (t *Tree) ToString() string {
	return fmt.Sprintf("%o %s %s %s\n", t.mode, t.GetType(), GetSha1AsString(t.GetHash()), t.name)
}

func (t *Tree) GetType() string {
	return "tree"
}

func (t *Tree) Children() []BGitObject {
	return t.entries
}

func (t *Tree) FormatToIndex() string {
	var buffer bytes.Buffer
	for _, entry := range t.entries {
		buffer.WriteString(entry.FormatToIndex())
	}
	return buffer.String()
}
