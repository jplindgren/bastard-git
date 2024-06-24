package object

type BGitObject interface {
	Serialize() []byte
	GetHash() []byte
	ToString() string
	GetType() string
	FormatToIndex() string
}
