package mqMessageModel

import "time"

type CollectNotes struct {
	Action    int
	Nid       string
	Uid       int64
	Timestamp time.Time
}
