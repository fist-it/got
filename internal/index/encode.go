package index

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

const (
	signature  = "DIRC"
	headerSize = 12
)

func (i *Index) Encode() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(signature)
	binary.Write(&buf, binary.BigEndian, i.Version)
	binary.Write(&buf, binary.BigEndian, uint32(len(i.Entries)))

	for _, e := range i.Entries {
		if err := encodeEntry(&buf, e); err != nil {
			return nil, err
		}
	}

	sum := sha256.Sum256(buf.Bytes())
	buf.Write(sum[:])

	return buf.Bytes(), nil
}

func encodeEntry(buf *bytes.Buffer, e *Entry) error {
	if len(e.Path) > 0xFFF {
		// path len encoded in low 12 bits of flags, use 0xFFF as overflow marker
	}

	start := buf.Len()

	binary.Write(buf, binary.BigEndian, uint32(e.CTime.Unix()))
	binary.Write(buf, binary.BigEndian, uint32(e.CTime.Nanosecond()))
	binary.Write(buf, binary.BigEndian, uint32(e.MTime.Unix()))
	binary.Write(buf, binary.BigEndian, uint32(e.MTime.Nanosecond()))
	binary.Write(buf, binary.BigEndian, e.Dev)
	binary.Write(buf, binary.BigEndian, e.Ino)
	binary.Write(buf, binary.BigEndian, e.Mode)
	binary.Write(buf, binary.BigEndian, e.Uid)
	binary.Write(buf, binary.BigEndian, e.Gid)
	binary.Write(buf, binary.BigEndian, e.Size)
	buf.Write(e.Hash[:])

	nameLen := uint16(len(e.Path))
	if len(e.Path) >= 0xFFF {
		nameLen = 0xFFF
	}
	flags := (e.Flags &^ 0xFFF) | nameLen
	binary.Write(buf, binary.BigEndian, flags)

	buf.WriteString(e.Path)
	buf.WriteByte(0)

	written := buf.Len() - start
	pad := (8 - written%8) % 8
	for range pad {
		buf.WriteByte(0)
	}

	return nil
}

func Decode(data []byte) (*Index, error) {
	if len(data) < headerSize+sha256.Size {
		return nil, fmt.Errorf("index: file too small")
	}

	body := data[:len(data)-sha256.Size]
	trailer := data[len(data)-sha256.Size:]

	sum := sha256.Sum256(body)
	if !bytes.Equal(sum[:], trailer) {
		return nil, fmt.Errorf("index: checksum mismatch")
	}

	if string(body[0:4]) != signature {
		return nil, fmt.Errorf("index: bad signature")
	}

	version := binary.BigEndian.Uint32(body[4:8])
	count := binary.BigEndian.Uint32(body[8:12])

	idx := &Index{Version: version, Entries: make([]*Entry, 0, count)}

	off := headerSize
	for range count {
		e, n, err := decodeEntry(body[off:])
		if err != nil {
			return nil, err
		}
		idx.Entries = append(idx.Entries, e)
		off += n
	}

	return idx, nil
}

func decodeEntry(data []byte) (*Entry, int, error) {
	const fixed = 4*10 + 32 + 2
	if len(data) < fixed {
		return nil, 0, fmt.Errorf("index: entry truncated")
	}

	// there has to be a better way to do it
	e := &Entry{}
	ctimeSec := binary.BigEndian.Uint32(data[0:4])
	ctimeNano := binary.BigEndian.Uint32(data[4:8])
	mtimeSec := binary.BigEndian.Uint32(data[8:12])
	mtimeNano := binary.BigEndian.Uint32(data[12:16])
	e.CTime = unixTime(ctimeSec, ctimeNano)
	e.MTime = unixTime(mtimeSec, mtimeNano)
	e.Dev = binary.BigEndian.Uint32(data[16:20])
	e.Ino = binary.BigEndian.Uint32(data[20:24])
	e.Mode = binary.BigEndian.Uint32(data[24:28])
	e.Uid = binary.BigEndian.Uint32(data[28:32])
	e.Gid = binary.BigEndian.Uint32(data[32:36])
	e.Size = binary.BigEndian.Uint32(data[36:40])
	copy(e.Hash[:], data[40:72])
	e.Flags = binary.BigEndian.Uint16(data[72:74])

	nameLen := int(e.Flags & 0xFFF)
	pathStart := fixed

	var path []byte
	if nameLen < 0xFFF {
		if len(data) < pathStart+nameLen+1 {
			return nil, 0, fmt.Errorf("index: path truncated")
		}
		path = data[pathStart : pathStart+nameLen]
	} else {
		end := bytes.IndexByte(data[pathStart:], 0)
		if end < 0 {
			return nil, 0, fmt.Errorf("index: unterminated path")
		}
		path = data[pathStart : pathStart+end]
	}
	e.Path = string(path)

	consumed := pathStart + len(path) + 1
	pad := (8 - consumed%8) % 8
	consumed += pad

	return e, consumed, nil
}
