package mqMessageModel

import "time"

type LikeNotes struct {
	Action    string
	Nid       string
	Uid       int
	Timestamp time.Time
}
