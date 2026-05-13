package index

import "time"

type Entry struct {
	CTime time.Time
	MTime time.Time
	Dev   uint32
	Ino   uint32
	Mode  uint32
	Uid   uint32
	Gid   uint32
	Size  uint32
	Hash  [32]byte
	Flags uint16
	Path  string
}
