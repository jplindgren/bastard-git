package repository

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type IndexEntry struct {
	name    string
	size    int64
	modTime string
}

func (r *Repository) addFileToIndex(filePath string) error {
	fmt.Println("add file to index" + filePath)
	return nil
}

func (r *Repository) readIndex() (map[string]IndexEntry, error) {
	indexEntries := make(map[string]IndexEntry)

	file, err := os.Open(r.paths.IndexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return indexEntries, nil
		}
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fLine := scanner.Text()
		data := strings.Split(fLine, " ")
		size, err := strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			return nil, err
		}

		modTime, err := time.Parse(time.RFC3339, data[3])
		if err != nil {
			return nil, err
		}

		indexEntries[data[2]] = IndexEntry{
			size:    size,
			name:    data[2],
			modTime: modTime.Format(time.RFC3339),
		}
	}

	return indexEntries, nil
}

func (r *Repository) updateIndex(buffer *bytes.Buffer) error {
	return os.WriteFile(r.paths.IndexPath, buffer.Bytes(), 0644)
}
