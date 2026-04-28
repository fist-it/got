package porcelain

import (
	"crypto/sha256"
	"os"

	"github.com/fist-it/got/internal/object"
	"github.com/fist-it/got/internal/store"
)

func HashObject(path string, write bool) ([32]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return [32]byte{}, err
	}

	blob := &object.Blob{Content: data}

	if write {
		s := &store.Store{Root: ".got/objects"}
		return s.Write(blob)
	}

	serialized := blob.Serialize()
	return sha256.Sum256(serialized), nil
}
