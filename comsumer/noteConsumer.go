package comsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"note_app_server_mq/config"
	"note_app_server_mq/config/action"
	"note_app_server_mq/global"
	"note_app_server_mq/model/mqMessageModel"
	"note_app_server_mq/repository"
	"note_app_server_mq/service"
	"sync"
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
	ctx := context.Background()
	conn, err := kafka.DialLeader(ctx, config.AC.Kafka.Network, config.AC.Kafka.Host+":"+config.AC.Kafka.Port, config.AC.Kafka.NoteLikes.Topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	for {
		msg := mqMessageModel.LikeNotes{}
		message, err1 := conn.ReadMessage(10e3)
		if err1 != nil {
			log.Fatal("failed to read message:", err1)
		}

		err = json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Fatal("failed to unmarshal message:", err)
		}

		fmt.Println(msg)

		uid := msg.Uid
		nid := msg.Nid
		if msg.Action == action.LikeNote {
			err = repository.LikeNote(nid, uid)
			if err != nil {
				log.Println("failed to like note:", err)
			}
		} else if msg.Action == action.DislikeNote {
			err = repository.CancelLikeNote(nid, uid)
			if err != nil {
				log.Println("failed to dislike note:", err)
			}
		}
	}
}

// 收藏&取消收藏笔记
func collectNote() {
	ctx := context.Background()
	conn, err := kafka.DialLeader(ctx, config.AC.Kafka.Network, config.AC.Kafka.Host+":"+config.AC.Kafka.Port, config.AC.Kafka.NoteCollects.Topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	for {
		msg := mqMessageModel.CollectNotes{}
		message, err1 := conn.ReadMessage(10e3)
		if err1 != nil {
			log.Println("failed to read message:", err1)
		}

		err = json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Println("failed to unmarshal message:", err)
		}

		fmt.Println(msg)

		uid := msg.Uid
		nid := msg.Nid
		if msg.Action == action.CollectNote {
			err = repository.CollectNote(nid, uid)
			if err != nil {
				log.Println("failed to collect note:", err)
			}
		} else if msg.Action == action.AbandonNote {
			err = repository.CancelCollectNote(nid, uid)
			if err != nil {
				log.Println("failed to Abandon note:", err)
			}
		}
	}
}

// 同步笔记
func syncNoteToES() {
	ctx := context.Background()
	conn, err := kafka.DialLeader(ctx, config.AC.Kafka.Network, config.AC.Kafka.Host+":"+config.AC.Kafka.Port, config.AC.Kafka.SyncNotes.Topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	for {
		msg := mqMessageModel.SyncNoteMsg{}
		message, err1 := conn.ReadMessage(10e4)
		if err1 != nil {
			log.Println("failed to read message:", err1)
		}

		err = json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Println("failed to unmarshal message:", err)
		}
		fmt.Println(msg)

		if msg.Action == action.SyncNote {
			err = repository.SaveNoteToES(msg.Note)
			if err != nil {
				// 待加入重试机制
				log.Println("failed to save note:", err)
			}
		}
	}
}

// 删除笔记
func delNote() {
	ctx := context.Background()
	conn, err := kafka.DialLeader(ctx, config.AC.Kafka.Network, config.AC.Kafka.Host+":"+config.AC.Kafka.Port, config.AC.Kafka.DelNotes.Topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	for {
		msg := mqMessageModel.DelNote{}
		message, err1 := conn.ReadMessage(10e3)
		if err1 != nil {
			log.Println("failed to read message:", err1)
		}

		err = json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Println("failed to unmarshal message:", err)
		}

		fmt.Println(msg)

		uid := msg.Uid
		nid := msg.Nid
		if msg.Action == action.DelNote {
			go func() {
				err = repository.DeleteNoteWithUid(nid, uid)
				if err != nil {
					log.Println("failed to collect note:", err)
				}
			}()
			go func() {
				_, err = global.ESClient.Delete("notes", nid).Do(ctx)
				if err != nil {
					log.Println("failed to delete ES data:", err)
				}
			}()
			go func() {
				err = service.DeleteDir(config.AC.Oss.NotePicsBucket, nid)
				if err != nil {
					log.Println("failed to delete dir:", err)
				}
			}()
		}
	}
}
