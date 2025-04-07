package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"note_app_server_mq/global"
	"note_app_server_mq/model/noteModel"
	"note_app_server_mq/model/userModel"
)

// LikeNote 点赞
func LikeNote(nid string, uid int) error {
	if err := global.Db.Where("uid = ? and nid = ?", uid, nid).First(&noteModel.LikedNotes{}).Error; err == nil {
		return errors.New("has liked")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	tx := global.Db.Begin()
	if err := tx.Create(&noteModel.LikedNotes{Uid: uid, Nid: nid}).Error; err != nil {
		tx.Rollback()
		return err
	}
	result := tx.Model(&noteModel.NoteInfo{}).Where("nid = ?", nid).UpdateColumn("likes_count", gorm.Expr("likes_count + ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("更新数据失败")
	}

	result = tx.Model(&userModel.UserCreationInfo{}).Where("uid = ?", uid).UpdateColumn("likes", gorm.Expr("likes + ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("更新数据失败")
	}
	tx.Commit()
	return nil
}

// CancelLikeNote 取消点赞
func CancelLikeNote(nid string, uid int) error {
	if err := global.Db.Model(&noteModel.LikedNotes{}).Where("uid = ? and nid = ?", uid, nid).First(&noteModel.LikedNotes{}).Error; err != nil {
		return errors.New("hasn't liked")
	}
	tx := global.Db.Begin()
	result := tx.Model(&noteModel.LikedNotes{}).Where("uid = ? and nid = ?", uid, nid).Delete(&noteModel.LikedNotes{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("取消点赞失败")
	}
	result = tx.Model(&noteModel.NoteInfo{}).Where("nid = ?", nid).UpdateColumn("likes_count", gorm.Expr("likes_count - ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("更新数据失败")
	}

	result = tx.Model(&userModel.UserCreationInfo{}).Where("uid = ?", uid).UpdateColumn("likes", gorm.Expr("likes - ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("更新数据失败")
	}
	tx.Commit()
	return nil
}

// GetNoteLikes 查询点赞数
func GetNoteLikes(nid string) (int, error) {
	noteInfo := noteModel.NoteInfo{}
	if err := global.Db.Where("nid = ?", nid).First(&noteInfo).Error; err != nil {
		return 0, err
	}
	return noteInfo.LikesCount, nil
}

// GetNoteCollections 查询收藏数
func GetNoteCollections(nid string) (int, error) {
	noteInfo := noteModel.NoteInfo{}
	if err := global.Db.Where("nid = ?", nid).First(&noteInfo).Error; err != nil {
		return 0, err
	}
	return noteInfo.CollectionsCount, nil
}

// CollectNote 收藏
func CollectNote(nid string, uid int) error {
	if err := global.Db.Model(&noteModel.CollectedNotes{}).Where("uid = ? and nid = ?", uid, nid).First(&noteModel.CollectedNotes{}).Error; err == nil {
		return errors.New("has collected")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	tx := global.Db.Begin()
	if err := tx.Create(&noteModel.CollectedNotes{Uid: uid, Nid: nid}).Error; err != nil {
		tx.Rollback()
		return err
	}
	result := tx.Model(&noteModel.NoteInfo{}).Where("nid = ?", nid).UpdateColumn("collections_count", gorm.Expr("collections_count + ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("更新数据失败")
	}

	result = tx.Model(&userModel.UserCreationInfo{}).Where("uid = ?", uid).UpdateColumn("collects", gorm.Expr("collects + ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("更新数据失败")
	}
	tx.Commit()
	return nil
}

// CancelCollectNote 取消收藏
func CancelCollectNote(nid string, uid int) error {
	if err := global.Db.Model(&noteModel.CollectedNotes{}).Where("uid = ? and nid = ?", uid, nid).First(&noteModel.CollectedNotes{}).Error; err != nil {
		return errors.New("hasn't liked")
	}
	tx := global.Db.Begin()
	result := tx.Model(&noteModel.CollectedNotes{}).Where("uid = ? and nid = ?", uid, nid).Delete(&noteModel.CollectedNotes{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("cancel collect failed")
	}

	result = tx.Model(&noteModel.NoteInfo{}).Where("nid = ?", nid).UpdateColumn("collections_count", gorm.Expr("collections_count - ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("更新数据失败")
	}

	result = tx.Model(&userModel.UserCreationInfo{}).Where("uid = ?", uid).UpdateColumn("collects", gorm.Expr("collects - ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("更新数据失败")
	}

	tx.Commit()
	return nil
}

// DeleteNoteWithUid 删除笔记
func DeleteNoteWithUid(nid string, uid int) error {
	result := global.Db.Where("nid = ? and uid = ?", nid, uid).Delete(&noteModel.Note{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("无匹配字段")
	}
	return nil
}

// LoadLikedNoteToRdb 将DB内的笔记点赞数据载入redis
func LoadLikedNoteToRdb(uid int) ([]noteModel.LikedNotes, error) {
	var likedNoteList []noteModel.LikedNotes
	if err := global.Db.Where("uid = ?", uid).Find(&likedNoteList).Error; err != nil {
		return nil, err
	}
	return likedNoteList, nil
}

// LoadCollectedNoteToRdb 将DB内的笔记收藏数据载入redis
func LoadCollectedNoteToRdb(uid int) ([]noteModel.CollectedNotes, error) {
	var collectedNoteList []noteModel.CollectedNotes
	if err := global.Db.Where("uid = ?", uid).Find(&collectedNoteList).Error; err != nil {
		return nil, err
	}
	return collectedNoteList, nil
}

// SaveNoteToES 转存笔记
func SaveNoteToES(note *noteModel.ESNote) error {
	_, err := global.ESClient.Create("notes", note.Nid).Request(note).Do(context.TODO())
	if err != nil {
		return err
	}
	return nil
}
