package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"note_app_server_mq/global"
	"note_app_server_mq/model/commentModel"
	"note_app_server_mq/model/noteModel"
)

// DeleteComment 删除评论
func DeleteComment(ctx context.Context, uid int64, cid string) error {
	cmt := &commentModel.Comment{}
	if err := global.Db.WithContext(ctx).Where("cid = ? and uid = ?", cid, uid).First(&cmt).Error; err != nil {
		return err
	}
	tx := global.Db.WithContext(ctx).Begin()
	result := tx.WithContext(ctx).Model(&commentModel.Comment{}).Where("cid = ? and uid = ?", cid, uid).Delete(&commentModel.Comment{})
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("无匹配记录")
	}
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	result = tx.WithContext(ctx).Model(&commentModel.Comment{}).Where("root_id = ?", cid).Delete(&commentModel.Comment{})
	if result.Error != nil {
		log.Println("删除子评论失败")
	}

	result = tx.WithContext(ctx).Model(&noteModel.NoteInfo{}).Where("nid = ?", cmt.Nid).UpdateColumn("comments_count", gorm.Expr("comments_count - ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("更新数据失败")
	}
	tx.WithContext(ctx).Commit()
	return nil
}

// LikeComment 点赞评论
func LikeComment(ctx context.Context, uid int64, cid string) error {
	if err := global.Db.WithContext(ctx).Where("uid = ? and cid = ?", uid, cid).First(&commentModel.LikedComment{}).Error; err == nil {
		return errors.New("已点赞")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	tx := global.Db.WithContext(ctx).Begin()
	result := tx.WithContext(ctx).Create(&commentModel.LikedComment{Uid: uid, Cid: cid})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.WithContext(ctx).Model(&commentModel.CommentsInfo{}).Where("cid = ?", cid).UpdateColumn("likes_count", gorm.Expr("likes_count + ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("更新数据失败")
	}
	tx.WithContext(ctx).Commit()
	return nil
}

// DislikeComment 取消点赞评论
func DislikeComment(ctx context.Context, uid int64, cid string) error {
	if err := global.Db.WithContext(ctx).Where("uid = ? and cid = ?", uid, cid).First(&commentModel.LikedComment{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("未点赞")
	} else if err != nil {
		return err
	}

	tx := global.Db.WithContext(ctx).Begin()
	result := tx.WithContext(ctx).Where("uid = ? and cid = ?", uid, cid).Delete(&commentModel.LikedComment{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("取消点赞失败")
	}

	result = tx.WithContext(ctx).Model(&commentModel.CommentsInfo{}).Where("cid = ?", cid).UpdateColumn("likes_count", gorm.Expr("likes_count - ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("更新数据失败")
	}
	tx.WithContext(ctx).Commit()
	return nil
}

// GetCommentLikes 获取评论点赞数
func GetCommentLikes(ctx context.Context, cid string) (int64, error) {
	comment := &commentModel.CommentsInfo{}
	if err := global.Db.WithContext(ctx).Model(&commentModel.CommentsInfo{}).Where("cid = ?", cid).First(&comment).Error; err != nil {
		return 0, err
	}
	return comment.LikesCount, nil
}

// GetUserLikedComment 获取用户点赞过的评论
func GetUserLikedComment(ctx context.Context, uid int64) ([]commentModel.LikedComment, error) {
	var likedComment []commentModel.LikedComment
	if err := global.Db.WithContext(ctx).Model(&commentModel.LikedComment{}).Where("uid = ?", uid).Scan(&likedComment).Error; err != nil {
		return nil, err
	}
	return likedComment, nil
}
