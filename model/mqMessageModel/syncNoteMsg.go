package mqMessageModel

import (
	"encoding/json"
	"note_app_server_mq/model/noteModel"
	"time"
)

type SyncNoteMsg struct {
	Action    int
	Note      *noteModel.ESNote
	Timestamp time.Time
}

func (that *SyncNoteMsg) EncodeMsg() ([]byte, error) {
	return json.Marshal(that)
}

func (that *SyncNoteMsg) DecodeMsg(bts []byte) error {
	return json.Unmarshal(bts, that)
}
