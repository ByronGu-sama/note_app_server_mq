package service

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"note_app_server_mq/global"
	"note_app_server_mq/repository"
	"strconv"
	"time"
)

// IncrNoteThumbsUp 增加笔记点赞数
func IncrNoteThumbsUp(ctx context.Context, uid string, nid string) {
	// 笔记点赞数
	thumbsUpNid := nid + ":ThumbsUp"
	// 点赞过的笔记
	userLikedNote := uid + ":Liked"

	var err error

	for i := 0; i < 3; i++ {
		_, err = global.NoteNormalRdb.Exists(ctx, thumbsUpNid).Result()
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

		global.NoteNormalRdb.Set(ctx, thumbsUpNid, dbLikesCnt, 0)
	}

	// 缓存中没有用户点赞过的笔记列表则从DB加载
	if global.NoteNormalRdb.SCard(ctx, userLikedNote).Val() == 0 {
		newUid, _ := strconv.Atoi(uid)
		ans, _ := repository.LoadLikedNoteToRdb(uint(newUid))
		for _, x := range ans {
			_, err = global.NoteNormalRdb.SAdd(ctx, userLikedNote, x.Nid).Result()
		}
	}

	// 先判断该用户是否已点赞过该笔记
	exist, err1 := global.NoteNormalRdb.SIsMember(ctx, userLikedNote, nid).Result()
	if err1 != nil {
		log.Println("Failed to get user thumbs up from rdb:", err1)
		return
	}

	if !exist {
		// 缓存仅在当天生效
		now := time.Now()
		midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
		duration := midnight.Sub(now)
		// 当前用户点赞列表中新增笔记
		_, err = global.NoteNormalRdb.SAdd(ctx, userLikedNote, nid).Result()

		newUid, _ := strconv.Atoi(uid)

		// 后台协程将新增的数据异步写入DB
		go func() {
			err := repository.LikeNote(nid, uint(newUid))
			for _ = range 3 {
				err = repository.LikeNote(nid, uint(newUid))
				if err == nil {
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
		}()

		// 笔记点赞数自增
		_, err = global.NoteNormalRdb.Incr(ctx, thumbsUpNid).Result()
		if err != nil {
			log.Println("Failed to increment thumbs up:", err)
		}

		global.NoteNormalRdb.Expire(ctx, userLikedNote, duration)
		if err != nil {
			log.Println("Failed to add user liked note:", err1)
		}
	}
}

// DecrNoteThumbsUp 减少笔记点赞数
func DecrNoteThumbsUp(ctx context.Context, uid string, nid string) {
	// 笔记点赞数
	thumbsUpNid := nid + ":ThumbsUp"
	// 点赞过的笔记
	userLikedNote := uid + ":Liked"
	var err error

	for i := 0; i < 3; i++ {
		_, err = global.NoteNormalRdb.Exists(ctx, thumbsUpNid).Result()
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

		global.NoteNormalRdb.Set(ctx, thumbsUpNid, dbLikesCnt, 0)
	}

	// 先判断该用户是否已点赞过该笔记
	exist, err1 := global.NoteNormalRdb.SIsMember(ctx, userLikedNote, nid).Result()
	if err1 != nil {
		log.Println("Failed to get user thumbs up from rdb:", err1)
		return
	}
	if exist {
		// 当前用户点赞列表中删除该笔记
		_, err = global.NoteNormalRdb.SRem(ctx, userLikedNote, nid).Result()
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

		newUid, _ := strconv.Atoi(uid)

		// 后台协程将减少的数据异步写入DB
		go func() {
			err := repository.CancelLikeNote(nid, uint(newUid))
			for _ = range 3 {
				err = repository.CancelLikeNote(nid, uint(newUid))
				if err == nil {
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
		}()

		keys := []string{userLikedNote}
		_, err = decrThumbsUpLuaScript.Run(ctx, global.NoteNormalRdb, keys).Result()
		if err != nil {
			log.Println("Failed to increment thumbs up:", err)
		}
	}
}

// IncrNoteCollection 增加笔记收藏数
func IncrNoteCollection(ctx context.Context, uid string, nid string) {
	// 笔记收藏数
	collectedNid := nid + ":Collection"
	// 收藏过的笔记
	userCollectedNote := uid + ":Collected"

	var err error

	// 查看该笔记收藏数是否在redis中已存在
	for i := 0; i < 3; i++ {
		_, err = global.NoteNormalRdb.Exists(ctx, collectedNid).Result()
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	// 如果重试后依然无法查询到redis中笔记对应的收藏数
	if err != nil {
		dbCollectionsCnt, e2 := repository.GetNoteCollections(nid)
		// 查数据库也失败时放弃减少收藏数据
		if e2 != nil {
			log.Fatal("Failed to get note collections from db:", e2)
			return
		}
		global.NoteNormalRdb.Set(ctx, collectedNid, dbCollectionsCnt, 0)
	}

	// 如果该值已被清空则从DB重新加载
	if global.NoteNormalRdb.SCard(ctx, userCollectedNote).Val() == 0 {
		newUid, _ := strconv.Atoi(uid)
		ans, _ := repository.LoadCollectedNoteToRdb(uint(newUid))
		for _, x := range ans {
			_, err = global.NoteNormalRdb.SAdd(ctx, userCollectedNote, x.Nid).Result()
		}
	}

	// 先判断该用户是否已收藏过该笔记
	exist, err1 := global.NoteNormalRdb.SIsMember(ctx, userCollectedNote, nid).Result()
	if err1 != nil {
		log.Println("Failed to get user collections from rdb:", err1)
		return
	}

	// 未收藏
	if !exist {
		// 缓存仅在当天生效
		now := time.Now()
		midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
		duration := midnight.Sub(now)
		// 当前用户点赞列表中增加该笔记
		_, err = global.NoteNormalRdb.SAdd(ctx, userCollectedNote, nid).Result()

		newUid, _ := strconv.Atoi(uid)

		// 后台协程将增加的数据异步写入DB
		go func() {
			err := repository.CollectNote(nid, uint(newUid))
			for _ = range 3 {
				err = repository.CollectNote(nid, uint(newUid))
				if err == nil {
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
		}()

		// 笔记点赞数自增
		_, err = global.NoteNormalRdb.Incr(ctx, collectedNid).Result()
		if err != nil {
			log.Println("Failed to increment collections:", err)
		}

		global.NoteNormalRdb.Expire(ctx, collectedNid, duration)
		if err != nil {
			log.Println("Failed add note to user collected rdb:", err1)
		}
	}
}

// DecrNoteCollection 减少笔记收藏数
func DecrNoteCollection(ctx context.Context, uid string, nid string) {
	// 笔记收藏数
	collectedNid := nid + ":Collection"
	// 收藏过的笔记
	userCollectedNote := uid + ":Collected"

	var err error

	for i := 0; i < 3; i++ {
		_, err = global.NoteNormalRdb.Exists(ctx, collectedNid).Result()
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
		global.NoteNormalRdb.Set(ctx, collectedNid, dbCollectionsCnt, 0)
	}

	// 先判断该用户是否已收藏过该笔记
	exist, err1 := global.NoteNormalRdb.SIsMember(ctx, userCollectedNote, nid).Result()
	if err1 != nil {
		log.Println("Failed to get user collections from rdb:", err1)
		return
	}
	if exist {
		// 当前用户点赞列表中删除该笔记
		_, err = global.NoteNormalRdb.SRem(ctx, userCollectedNote, nid).Result()
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

		newUid, _ := strconv.Atoi(uid)

		// 后台协程将减少的数据异步写入DB
		go func() {
			err := repository.CancelCollectNote(nid, uint(newUid))
			for _ = range 3 {
				err = repository.CancelCollectNote(nid, uint(newUid))
				if err == nil {
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
		}()

		keys := []string{collectedNid}
		_, err = decrCollectionsLuaScript.Run(ctx, global.NoteNormalRdb, keys).Result()
		if err != nil {
			log.Println("Failed to decrement thumbs up:", err)
		}
	}
}
