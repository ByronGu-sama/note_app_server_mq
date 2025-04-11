package comsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"note_app_server_mq/config"
	"note_app_server_mq/config/action"
	"note_app_server_mq/model/mqMessageModel"
	"note_app_server_mq/repository"
	"note_app_server_mq/service"
	"note_app_server_mq/utils"
	"strconv"
	"sync"
	"time"
)

func NoteListener() {
	fmt.Println("note MQ listener has started")
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		likeNote()
	}()
	go func() {
		defer wg.Done()
		collectNote()
	}()
	go func() {
		defer wg.Done()
		syncNoteToES()
	}()
	go func() {
		defer wg.Done()
		delNote()
	}()
	wg.Wait()
}

// LikeNote 点赞&取消点赞笔记
func likeNote() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.AC.Kafka.Host + ":" + config.AC.Kafka.Port},
		GroupID:  "note_thumbsUp_group",
		Topic:    config.AC.Kafka.NoteLikes.Topic,
		MaxBytes: 10e3,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		msg := mqMessageModel.LikeNotes{}
		message, err1 := reader.ReadMessage(ctx)
		if err1 != nil {
			log.Fatal("failed to read message:", err1)
		}
		err := json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Fatal("failed to unmarshal message:", err)
		}

		log.Println(fmt.Sprintf("Fetched New Msg:%v", msg))

		uid := msg.Uid
		nid := msg.Nid

		act, _ := strconv.Atoi(msg.Action)
		switch act {
		case action.LikeNote:
			go func(uid int64, nid string) {
				utils.SafeGo(func() {
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					defer cancel()
					// 更新缓存中的点赞数
					service.IncrNoteThumbsUp(ctx, uid, nid)
				})
			}(uid, nid)
		case action.DislikeNote:
			go func(uid int64, nid string) {
				utils.SafeGo(func() {
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					defer cancel()
					// 处理取消点赞帖子的逻辑
					service.DecrNoteThumbsUp(ctx, uid, nid)
				})
			}(uid, nid)
		default:
			log.Println("not pre defined case")
		}
	}
}

// 收藏&取消收藏笔记
func collectNote() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.AC.Kafka.Host + ":" + config.AC.Kafka.Port},
		GroupID:  "note_collection_group",
		Topic:    config.AC.Kafka.NoteCollects.Topic,
		MaxBytes: 10e3,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		message, err1 := reader.ReadMessage(ctx)

		if err1 != nil {
			log.Println("failed to read message:", err1)
		}

		msg := mqMessageModel.CollectNotes{}
		err := json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Println("failed to unmarshal message:", err)
		}

		log.Println(fmt.Sprintf("Fetched New Msg:%v", msg))

		uid := msg.Uid
		nid := msg.Nid

		act, _ := strconv.Atoi(msg.Action)
		switch act {
		case action.CollectNote:
			go func(uid int64, nid string) {
				utils.SafeGo(func() {
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					defer cancel()
					service.IncrNoteCollection(ctx, uid, nid)
				})
			}(uid, nid)
		case action.AbandonNote:
			go func(uid int64, nid string) {
				utils.SafeGo(func() {
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					defer cancel()
					service.DecrNoteCollection(ctx, uid, nid)
				})
			}(uid, nid)
		default:
			log.Println("not pre defined case")
		}
	}
}

// 同步笔记
func syncNoteToES() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.AC.Kafka.Host + ":" + config.AC.Kafka.Port},
		GroupID:  "note_sync_note_group",
		Topic:    config.AC.Kafka.SyncNotes.Topic,
		MaxBytes: 10e6,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		message, err1 := reader.ReadMessage(ctx)
		if err1 != nil {
			log.Println("failed to read message:", err1)
		}

		msg := mqMessageModel.SyncNoteMsg{}
		err := json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Println("failed to unmarshal message:", err)
		}

		log.Println(fmt.Sprintf("Fetched New Msg:%v", msg))

		act, _ := strconv.Atoi(msg.Action)
		switch act {
		case action.SyncNote:
			err = repository.SaveNoteToES(ctx, msg.Note)
			if err != nil {
				// 待加入重试机制
				log.Println("failed to save note:", err)
			}
		default:
			log.Println("not pre defined case")
		}
	}
}

// 删除笔记
func delNote() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.AC.Kafka.Host + ":" + config.AC.Kafka.Port},
		GroupID:  "note_del_group",
		Topic:    config.AC.Kafka.DelNotes.Topic,
		MaxBytes: 10e3,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		message, err1 := reader.ReadMessage(ctx)
		if err1 != nil {
			log.Println("failed to read message:", err1)
		}

		msg := mqMessageModel.DelNote{}
		err := json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Println("failed to unmarshal message:", err)
		}

		log.Println(fmt.Sprintf("Fetched New Msg:%v", msg))

		act, _ := strconv.Atoi(msg.Action)
		switch act {
		case action.DelNote:
			go func() {
				uid := msg.Uid
				nid := msg.Nid
				utils.SafeGo(func() {
					service.DelNote(ctx, uid, nid)
				})
			}()
		default:
			log.Println("not pre defined case")
		}
	}
}
