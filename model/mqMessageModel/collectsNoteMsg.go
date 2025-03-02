package mqMessageModel

import "time"

type CollectNotes struct {
	Action    string
	Nid       string
	Uid       uint
	Timestamp time.Time
}
