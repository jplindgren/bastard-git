package object

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
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

var CommitPrefix = []byte("c\000")

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
	buffer.WriteString(string(CommitPrefix))
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

func (c *Commit) FormatToIndex() string {
	return ""
}

func ParseCommit(data []byte, commitHash string) (*Commit, error) {
	data, found := bytes.CutPrefix(data, CommitPrefix)
	if !found {
		return nil, errors.New("could not parse object as commit: " + commitHash)
	}

	commit := &Commit{}

	sData := string(data)
	lines := strings.Split(sData, "\n")

	commit.Hash = []byte(commitHash)
	commit.Tree = []byte(strings.TrimPrefix(lines[0], "tree "))
	commit.Parent = strings.Split(strings.TrimPrefix(lines[1], "parent "), " ")
	commit.Author = strings.Replace(lines[2], "author ", "", 1)
	commit.Message = lines[4]
	return commit, nil
}
