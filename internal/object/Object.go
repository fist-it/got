package object

type Object interface {
	Type() string
	Serialize() []byte
}
