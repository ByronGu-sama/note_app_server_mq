package commentModel

type LikedComment struct {
	Uid int    `json:"uid" gorm:"uid"`
	Cid string `json:"cid" gorm:"cid"`
}

func (LikedComment) TableName() string {
	return "liked_comment"
}
