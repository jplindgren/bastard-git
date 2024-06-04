package object

import (
	"fmt"
	"io/fs"
	"time"
)

type Blob struct {
	content   string
	name      string
	hash      []byte
	mode      fs.FileMode
	size      int64
	modTime   time.Time
	sytemPath string
}

func NewBlob(content string, name string, mode fs.FileMode, modTime time.Time, size int64, sytemPath string) *Blob {
	hash := generateSHA1Hash(content)

	blob := &Blob{
		hash:      hash,
		content:   content,
		name:      name,
		mode:      mode,
		modTime:   modTime,
		size:      size,
		sytemPath: sytemPath,
	}

	fmt.Printf("Creating new blob for: %s hash: %s \n", blob.name, GetSha1AsString(blob.hash))
	return blob
}

func (b *Blob) Serialize() []byte {
	return []byte(b.content)
}

func (b *Blob) GetHash() []byte {
	return b.hash
}

func (b *Blob) ToString() string {
	return fmt.Sprintf("%o %s %s %s\n", b.mode, b.GetType(), GetSha1AsString(b.GetHash()), b.name)
}

func (b *Blob) FormatToIndex() string {
	return fmt.Sprintf("%s %d %s %s\n", GetSha1AsString(b.GetHash()), b.size, b.sytemPath, b.modTime.Format(time.RFC3339))
}

func (b *Blob) GetType() string {
	return "blob"
}

func (b *Blob) Children() []BGitObject {
	return []BGitObject{}
}
