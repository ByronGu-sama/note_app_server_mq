package noteModel

type NoteInfo struct {
	Nid              string `json:"nid" gorm:"column:nid"`
	LikesCount       int    `json:"likesCount" gorm:"column:likes_count"`
	CommentsCount    int    `json:"commentsCount" gorm:"column:comments_count"`
	CollectionsCount int    `json:"collectionsCount" gorm:"column:collections_count"`
	SharesCount      int    `json:"sharesCount" gorm:"column:shares_count"`
	ViewsCount       int    `json:"viewsCount" gorm:"column:views_count"`
}

func (NoteInfo) TableName() string {
	return "notes_info"
}
