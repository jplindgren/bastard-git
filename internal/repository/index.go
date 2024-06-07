package repository

import (
	"bytes"
	"fmt"
	"os"
)

type indexEntry struct {
	mode         uint32
	hash         string
	relativePath string
}

func (r *Repository) addFileToIndex(filePath string) error {
	fmt.Println("add file to index" + filePath)
	return nil
}

func (r *Repository) readIndex() error {
	fmt.Println("reading index")
	return nil
}

// TODO: kill old approach
func (r *Repository) updateIndexNew(buffer *bytes.Buffer) error {
	return os.WriteFile(r.Paths.IndexPath, buffer.Bytes(), 0644)
}
