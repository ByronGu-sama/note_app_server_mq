package mqMessageModel

import "time"

type LikeNotes struct {
	Action    int
	Nid       string
	Uid       int64
	Timestamp time.Time
}
