package mqMessageModel

import "time"

type CollectNotes struct {
	Action    string
	Nid       string
	Uid       int
	Timestamp time.Time
}
