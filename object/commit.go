package object

import (
	"bytes"
	"fmt"
	"time"
)

type Commit struct {
	hash      []byte
	author    string
	message   string
	createdAt time.Time
	tree      []byte   //sha1 of the root tree
	parent    []string //array of sha1 commits
}

func NewCommit(author, message string, tree []byte, parents []string) *Commit {
	createdAt := time.Now()
	hash := generateSHA1Hash(message + author + createdAt.String())

	commit := &Commit{
		author:    author,
		hash:      hash,
		message:   message,
		createdAt: createdAt,
		tree:      tree,
		parent:    parents,
	}
	fmt.Printf("Creating new commit with message %s at %s \n", commit.message, commit.createdAt.String())
	return commit
}

func (c *Commit) Serialize() []byte {
	var buffer bytes.Buffer
	//buffer.WriteString("c\000")
	buffer.WriteString(fmt.Sprintf("tree %s\n", GetSha1AsString(c.tree)))

	if len(c.parent) > 0 { //TODO: treat multiple parents
		buffer.WriteString(fmt.Sprintf("parent %s\n", c.parent[0]))
	}

	buffer.WriteString(fmt.Sprintf("author %s\n", c.author))
	buffer.WriteString(fmt.Sprintf("\n%s", c.message))

	return buffer.Bytes()
}

func (c *Commit) GetHash() []byte {
	return c.hash
}

func (c *Commit) ToString() string {
	return GetSha1AsString(c.hash)
}

func (c *Commit) GetType() string {
	return "commit"
}

func (c *Commit) Children() []BGitObject {
	return []BGitObject{}
}

func (c *Commit) FormatToIndex() string {
	return ""
}
