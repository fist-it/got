package object

import (
	"time"
)

type CommitList struct {
	C    *Commit
	Next *CommitList
}

type Commit struct {
	Timestamp time.Time
	Index     uint32
	Parents   *CommitList
}
