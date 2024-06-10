package object

import (
	"bytes"
	"fmt"
	"time"
)

type Commit struct {
	Hash      []byte
	Author    string
	Message   string
	CreatedAt time.Time
	Tree      []byte   //sha1 of the root tree
	Parent    []string //array of sha1 commits
}

func NewCommit(author, message string, tree []byte, parents []string) *Commit {
	createdAt := time.Now()
	hash := generateSHA1Hash(message + author + createdAt.String())

	commit := &Commit{
		Author:    author,
		Hash:      hash,
		Message:   message,
		CreatedAt: createdAt,
		Tree:      tree,
		Parent:    parents,
	}
	fmt.Printf("Creating new commit with message %s at %s \n", commit.Message, commit.CreatedAt.String())
	return commit
}

func (c *Commit) Serialize() []byte {
	var buffer bytes.Buffer
	//buffer.WriteString("c\000")
	buffer.WriteString(fmt.Sprintf("tree %s\n", GetSha1AsString(c.Tree)))

	if len(c.Parent) > 0 { //TODO: treat multiple parents
		buffer.WriteString(fmt.Sprintf("parent %s\n", c.Parent[0]))
	}

	buffer.WriteString(fmt.Sprintf("author %s\n", c.Author))
	buffer.WriteString(fmt.Sprintf("\n%s", c.Message))

	return buffer.Bytes()
}

func (c *Commit) GetHash() []byte {
	return c.Hash
}

func (c *Commit) ToString() string {
	return GetSha1AsString(c.Hash)
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
