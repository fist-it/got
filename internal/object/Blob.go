package object

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"strconv"
)

type Hash [32]byte

type Blob struct {
	Content []byte
}

func (b *Blob) Type() string {
	return "blob"
}

func (b *Blob) Serialize() []byte {

	var result []byte

	size_b := strconv.FormatInt(int64(len(b.Content)), 10)
	result = []byte("blob ")

	result = append(result, size_b...)
	result = append(result, 0)

	result = append(result, b.Content...)

	return result
}

func DeserializeBlob(data []byte) (*Blob, error) {
	header := data[0:5]

	if string(header) != "blob " {
		return nil, errors.New("Invalid blob header")
	}
	nullIdx := bytes.IndexByte(data, 0)
	sizeStr := string(data[5:nullIdx])

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return nil, errors.New("Invalid blob size")
	}

	content := data[nullIdx+1 : int64(nullIdx)+1+size]
	b := &Blob{
		Content: []byte(content),
	}

	return b, nil
}

func (b *Blob) Hash() Hash {
	return sha256.Sum256(b.Serialize())
}
