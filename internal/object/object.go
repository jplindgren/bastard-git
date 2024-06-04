package object

type BGitObject interface {
	Serialize() []byte
	GetHash() []byte
	ToString() string
	GetType() string
	Children() []BGitObject
	FormatToIndex() string
}
