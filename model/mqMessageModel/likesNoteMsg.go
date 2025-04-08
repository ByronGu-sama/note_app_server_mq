package mqMessageModel

import "time"

type LikeNotes struct {
	Action    string
	Nid       string
	Uid       int64
	Timestamp time.Time
}
