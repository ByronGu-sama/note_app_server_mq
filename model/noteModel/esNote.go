package noteModel

import (
	"encoding/json"
	"log"
	"time"
)

type Note struct {
	Nid         string    `json:"nid" gorm:"column:nid"`
	Uid         int64     `json:"uid" gorm:"column:uid"`
	Cover       string    `json:"cover" gorm:"column:cover"`
	CoverHeight float64   `json:"cover_height" gorm:"column:cover_height"`
	Pics        string    `json:"pics" gorm:"column:pics"`
	Title       string    `json:"title" gorm:"column:title"`
	Content     string    `json:"content" gorm:"column:content"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at"`
	Public      int       `json:"public" gorm:"column:public"`
	CategoryId  int       `json:"categoryId" gorm:"column:category_id"`
	Tags        string    `json:"tags" gorm:"column:tags"`
	Status      int       `json:"status" gorm:"column:status;default:1"`
}

func (Note) TableName() string {
	return "notes"
}

type ESNote struct {
	Nid         string    `json:"nid"`
	Uid         int64     `json:"uid"`
	Username    string    `json:"username"`
	AvatarUrl   string    `json:"avatarUrl"`
	Cover       string    `json:"cover"`
	CoverHeight float64   `json:"cover_height"`
	Pics        string    `json:"pics"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	LikesCount  int64     `json:"likes_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Public      bool      `json:"public"`
	CategoryId  int       `json:"category_id"`
	Tags        string    `json:"tags"`
	Status      int       `json:"status"`
}

func (that *ESNote) ToRawJson() []byte {
	result, err := json.Marshal(that)
	if err != nil {
		log.Fatal("jsonify failed")
	}
	return result
}

func (that *ESNote) ToJson() string {
	return string(that.ToRawJson())
}
