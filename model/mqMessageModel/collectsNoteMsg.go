package mqMessageModel

import "time"

type CollectNotes struct {
	Action    string
	Nid       string
	Uid       int64
	Timestamp time.Time
}
