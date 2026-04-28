package porcelain

import (
	"github.com/fist-it/got/internal/store"
)

func CatFile(hash [32]byte) ([]byte, error) {
	s := &store.Store{Root: ".got/objects"}
	return s.Read(hash)
}
