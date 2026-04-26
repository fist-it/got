package store

import (
	"compress/zlib"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/fist-it/got/internal/object"
)

type Store struct {
	Root string
}

func (s *Store) buildObjectPath(hash [32]byte) string {
	hex_hash := hex.EncodeToString(hash[:])

	path := filepath.Join(s.Root, hex_hash[:2])
	file_path := filepath.Join(path, hex_hash[2:])

	return file_path
}

func (s *Store) Write(o object.Object) ([32]byte, error) {
	data := o.Serialize()
	hash := sha256.Sum256(data)

	path := s.buildObjectPath(hash)
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return [32]byte{}, err
	}

	file, err := os.Create(path)
	if err != nil {
		return [32]byte{}, err
	}
	defer file.Close()

	zlibWriter := zlib.NewWriter(file)
	_, err = zlibWriter.Write(data)
	if err != nil {
		return [32]byte{}, nil
	}
	zlibWriter.Close()

	return hash, nil
}

func (s *Store) Read(hash [32]byte) ([]byte, error) {
	file_path := s.buildObjectPath(hash)

	file, err := os.Open(file_path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	zlibReader, err := zlib.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer zlibReader.Close()

	data, err := io.ReadAll(zlibReader)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Store) Exists(hash [32]byte) bool {
	file_path := s.buildObjectPath(hash)
	_, err := os.Stat(file_path)
	if err != nil {
		return false
	}
	return true
}
