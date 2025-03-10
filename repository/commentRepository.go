package repository

import (
	"errors"
	"gorm.io/gorm"
	"note_app_server_mq/global"
	"note_app_server_mq/model/commentModel"
	"note_app_server_mq/model/noteModel"
)

// DeleteComment 删除评论
func DeleteComment(uid uint, cid string) error {
	cmt := new(commentModel.Comment)
	if err := global.Db.Where("cid = ? and uid = ?", cid, uid).First(&cmt).Error; err != nil {
		return err
	}
	tx := global.Db.Begin()
	result := tx.Model(&commentModel.Comment{}).Where("cid = ? and uid = ?", cid, uid).Delete(&commentModel.Comment{})
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("无匹配记录")
	}
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Model(&noteModel.NoteInfo{}).Where("nid = ?", cmt.Nid).UpdateColumn("comments_count", gorm.Expr("comments_count - ?", 1))
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

// LikeComment 点赞评论
func LikeComment(uid uint, cid string) error {
	if err := global.Db.Where("uid = ? and cid = ?", uid, cid).First(&commentModel.LikedComment{}).Error; err == nil {
		return errors.New("已点赞")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	tx := global.Db.Begin()
	result := tx.Create(&commentModel.LikedComment{Uid: uid, Cid: cid})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Model(&commentModel.CommentsInfo{}).Where("cid = ?", cid).UpdateColumn("likes_count", gorm.Expr("likes_count + ?", 1))
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

// DislikeComment 取消点赞评论
func DislikeComment(uid uint, cid string) error {
	if err := global.Db.Where("uid = ? and cid = ?", uid, cid).First(&commentModel.LikedComment{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("未点赞")
	} else if err != nil {
		return err
	}

	tx := global.Db.Begin()
	result := tx.Where("uid = ? and cid = ?", uid, cid).Delete(&commentModel.LikedComment{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("取消点赞失败")
	}

	result = tx.Model(&commentModel.CommentsInfo{}).Where("cid = ?", cid).UpdateColumn("likes_count", gorm.Expr("likes_count - ?", 1))
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
