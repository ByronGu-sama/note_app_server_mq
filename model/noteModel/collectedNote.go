package noteModel

type CollectedNotes struct {
	Uid int    `json:"uid" gorm:"uid"`
	Nid string `json:"nid" gorm:"nid"`
}

func (CollectedNotes) TableName() string {
	return "collected_notes"
}
