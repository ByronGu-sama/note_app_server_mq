package service

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"note_app_server_mq/global"
	"note_app_server_mq/repository"
	"time"
)

// IncrNoteThumbsUp 增加笔记点赞数
func IncrNoteThumbsUp(ctx context.Context, uid string, nid string) {
	var err error

	for i := 0; i < 3; i++ {
		_, err = global.ThumbsUpRdbClient.Exists(ctx, nid).Result()
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if err != nil {
		dbLikesCnt, e2 := repository.GetNoteLikes(nid)
		// 查数据库也失败时放弃新增点赞数据
		if e2 != nil {
			log.Fatal("Failed to get note thumbs up from db:", e2)
			return
		}
		global.ThumbsUpRdbClient.Set(ctx, nid, dbLikesCnt, 0)
	}

	// 先判断该用户是否已点赞过该笔记
	exist, err1 := global.UserLikedNotesRdbClient.SIsMember(ctx, uid, nid).Result()
	if err1 != nil {
		log.Println("Failed to get user thumbs up from rdb:", err1)
		return
	}
	if !exist {
		// 当前用户点赞列表中新增笔记
		_, err = global.UserLikedNotesRdbClient.SAdd(ctx, uid, nid).Result()
		if err != nil {
			log.Println("Failed to add user liked note:", err1)
		}
		// 笔记点赞数自增
		_, err = global.ThumbsUpRdbClient.Incr(ctx, nid).Result()
		if err != nil {
			log.Println("Failed to increment thumbs up:", err)
		}
	}
}

// DecrNoteThumbsUp 减少笔记点赞数
func DecrNoteThumbsUp(ctx context.Context, uid string, nid string) {
	var err error

	for i := 0; i < 3; i++ {
		_, err = global.ThumbsUpRdbClient.Exists(ctx, nid).Result()
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if err != nil {
		dbLikesCnt, e2 := repository.GetNoteLikes(nid)
		// 查数据库也失败时放弃减少点赞数据
		if e2 != nil {
			log.Fatal("Failed to get note thumbs up from db:", e2)
			return
		}
		global.ThumbsUpRdbClient.Set(ctx, nid, dbLikesCnt, 0)
	}

	// 先判断该用户是否已点赞过该笔记
	exist, err1 := global.UserLikedNotesRdbClient.SIsMember(ctx, uid, nid).Result()
	if err1 != nil {
		log.Println("Failed to get user thumbs up from rdb:", err1)
		return
	}
	if exist {
		// 当前用户点赞列表中删除该笔记
		_, err = global.UserLikedNotesRdbClient.SRem(ctx, uid, nid).Result()
		if err != nil {
			log.Println("Failed to remove user liked note:", err1)
		}
		// 笔记点赞数减一
		decrThumbsUpLuaScript := redis.NewScript(`
			local cnt = redis.call('GET', KEYS[1])
			if cnt and tonumber(cnt) > 0 then
				if redis.call('DECR', KEYS[1]) == 0 then
					redis.call('DEL', KEYS[1])
				end
				return cnt
			end
			return 0
		`)

		keys := []string{nid}
		_, err = decrThumbsUpLuaScript.Run(ctx, global.ThumbsUpRdbClient, keys).Result()
		if err != nil {
			log.Println("Failed to increment thumbs up:", err)
		}
	}
}

// IncrNoteCollection 增加笔记收藏数
func IncrNoteCollection(ctx context.Context, uid string, nid string) {
	var err error

	for i := 0; i < 3; i++ {
		_, err = global.CollectedCntRdbClient.Exists(ctx, nid).Result()
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if err != nil {
		dbCollectionsCnt, e2 := repository.GetNoteCollections(nid)
		// 查数据库也失败时放弃减少收藏数据
		if e2 != nil {
			log.Fatal("Failed to get note collections from db:", e2)
			return
		}
		global.CollectedCntRdbClient.Set(ctx, nid, dbCollectionsCnt, 0)
	}

	// 先判断该用户是否已收藏过该笔记
	exist, err1 := global.UserCollectedNotesRdbClient.SIsMember(ctx, uid, nid).Result()
	if err1 != nil {
		log.Println("Failed to get user collections from rdb:", err1)
		return
	}
	if !exist {
		// 当前用户点赞列表中增加该笔记
		_, err = global.UserCollectedNotesRdbClient.SAdd(ctx, uid, nid).Result()
		if err != nil {
			log.Println("Failed add note to user collected rdb:", err1)
		}

		// 笔记点赞数自增
		_, err = global.CollectedCntRdbClient.Incr(ctx, nid).Result()
		if err != nil {
			log.Println("Failed to increment collections:", err)
		}
	}
}

// DecrNoteCollection 减少笔记收藏数
func DecrNoteCollection(ctx context.Context, uid string, nid string) {
	var err error

	for i := 0; i < 3; i++ {
		_, err = global.CollectedCntRdbClient.Exists(ctx, nid).Result()
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if err != nil {
		dbCollectionsCnt, e2 := repository.GetNoteCollections(nid)
		// 查数据库也失败时放弃减少收藏数据
		if e2 != nil {
			log.Fatal("Failed to get note collections from db:", e2)
			return
		}
		global.CollectedCntRdbClient.Set(ctx, nid, dbCollectionsCnt, 0)
	}

	// 先判断该用户是否已收藏过该笔记
	exist, err1 := global.UserCollectedNotesRdbClient.SIsMember(ctx, uid, nid).Result()
	if err1 != nil {
		log.Println("Failed to get user collections from rdb:", err1)
		return
	}
	if exist {
		// 当前用户点赞列表中删除该笔记
		_, err = global.UserCollectedNotesRdbClient.SRem(ctx, uid, nid).Result()
		if err != nil {
			log.Println("Failed add note to user collected rdb:", err1)
		}

		decrCollectionsLuaScript := redis.NewScript(`
			local cnt = redis.call('GET', KEYS[1])
			if cnt and tonumber(cnt) > 0 then
				if redis.call('DECR', KEYS[1]) == 0 then
					redis.call('DEL', KEYS[1])
				end
				return cnt
			end
			return 0
		`)

		keys := []string{nid}
		_, err = decrCollectionsLuaScript.Run(ctx, global.CollectedCntRdbClient, keys).Result()
		if err != nil {
			log.Println("Failed to decrement thumbs up:", err)
		}
	}
}
