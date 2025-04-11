package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"note_app_server_mq/global"
	"note_app_server_mq/repository"
	"note_app_server_mq/utils"
	"time"
)

// IncrCommentThumbsUp 增加评论点赞数
func IncrCommentThumbsUp(ctx context.Context, uid int64, cid string) {
	select {
	case <-ctx.Done():
		return
	default:
		// 评论点赞数
		thumbsUpCid := fmt.Sprintf("%s:ThumbsUp", cid)
		// 点赞过的评论
		userLikedComment := fmt.Sprintf("%d:Liked", uid)

		var err error

		for i := 0; i < 3; i++ {
			_, err = global.CommentNormalRdb.Exists(ctx, thumbsUpCid).Result()
			if err == nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}

		if err != nil {
			dbLikesCnt, e2 := repository.GetCommentLikes(ctx, cid)
			// 查数据库也失败时放弃新增点赞数据
			if e2 != nil {
				log.Fatal("Failed to get comment thumbs up from db:", e2)
				return
			}

			global.CommentNormalRdb.Set(ctx, thumbsUpCid, dbLikesCnt, 0)
		}

		// 缓存中没有用户点赞过的评论列表则从DB加载
		if global.CommentNormalRdb.SCard(ctx, userLikedComment).Val() == 0 {
			ans, _ := repository.GetUserLikedComment(ctx, uid)
			for _, x := range ans {
				_, err = global.CommentNormalRdb.SAdd(ctx, userLikedComment, x.Cid).Result()
			}
		}

		// 先判断该用户是否已点赞过该笔记
		exist, err1 := global.CommentNormalRdb.SIsMember(ctx, userLikedComment, cid).Result()
		if err1 != nil {
			log.Println("Failed to get user thumbs up from rdb:", err1)
			return
		}

		if !exist {
			// 缓存仅在当天生效
			now := time.Now()
			midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
			duration := midnight.Sub(now)
			// 当前用户点赞列表中新增评论
			_, err = global.CommentNormalRdb.SAdd(ctx, userLikedComment, cid).Result()

			// 后台协程将新增的数据异步写入DB
			go func() {
				utils.SafeGo(func() {
					e1 := repository.LikeComment(ctx, uid, cid)
					for _ = range 3 {
						e1 = repository.LikeComment(ctx, uid, cid)
						if e1 == nil {
							break
						}
						time.Sleep(500 * time.Millisecond)
					}
				})
			}()

			// 评论点赞数自增
			_, err = global.CommentNormalRdb.Incr(ctx, thumbsUpCid).Result()
			if err != nil {
				log.Println("Failed to increment thumbs up:", err)
			}

			global.CommentNormalRdb.Expire(ctx, userLikedComment, duration)
			if err != nil {
				log.Println("Failed to add user liked comment:", err1)
			}
		}
	}
}

// DecrCommentThumbsUp 减少评论点赞数
func DecrCommentThumbsUp(ctx context.Context, uid int64, cid string) {
	select {
	case <-ctx.Done():
		return
	default:
		// 评论点赞数
		thumbsUpCid := cid + ":ThumbsUp"
		// 点赞过的评论
		userLikedComment := fmt.Sprintf("%d:Liked", uid)

		var err error

		for i := 0; i < 3; i++ {
			_, err = global.CommentNormalRdb.Exists(ctx, thumbsUpCid).Result()
			if err == nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}

		if err != nil {
			dbLikesCnt, e2 := repository.GetCommentLikes(ctx, cid)
			// 查数据库也失败时放弃新增点赞数据
			if e2 != nil {
				log.Fatal("Failed to get comment thumbs up from db:", e2)
				return
			}

			global.CommentNormalRdb.Set(ctx, thumbsUpCid, dbLikesCnt, 0)
		}

		// 缓存中没有用户点赞过的评论列表则从DB加载
		if global.CommentNormalRdb.SCard(ctx, userLikedComment).Val() == 0 {
			ans, _ := repository.GetUserLikedComment(ctx, uid)
			for _, x := range ans {
				_, err = global.CommentNormalRdb.SAdd(ctx, userLikedComment, x.Cid).Result()
			}
		}

		// 先判断该用户是否已点赞过该笔记
		exist, err1 := global.CommentNormalRdb.SIsMember(ctx, userLikedComment, cid).Result()
		if err1 != nil {
			log.Println("Failed to get user thumbs up from rdb:", err1)
			return
		}

		if exist {
			// 当前用户点赞列表中删除该笔记
			_, err = global.CommentNormalRdb.SRem(ctx, userLikedComment, cid).Result()
			if err != nil {
				log.Println("Failed to remove user liked comment:", err1)
			}
			// 笔记点赞数减一
			decrThumbsUpLuaScript := redis.NewScript(`
			local cnt = redis.call('GET', KEYS[1])
			if cnt and tonumber(cnt) > 0 then
				local newCnt = redis.call('DECR', KEYS[1])
				if newCnt == 0 then
					redis.call('DEL', KEYS[1])
				end
				return newCnt
			end
			return 0
		`)

			// 后台协程将新增的数据异步写入DB
			go func() {
				utils.SafeGo(func() {
					e1 := repository.DislikeComment(ctx, uid, cid)
					for _ = range 3 {
						e1 = repository.DislikeComment(ctx, uid, cid)
						if e1 == nil {
							break
						}
						time.Sleep(500 * time.Millisecond)
					}
				})
			}()

			keys := []string{thumbsUpCid}
			_, err = decrThumbsUpLuaScript.Run(ctx, global.CommentNormalRdb, keys).Result()
			if err != nil {
				log.Println("Failed to decrement thumbs up:", err)
			}
		}
	}
}
