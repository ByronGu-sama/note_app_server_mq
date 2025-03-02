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
	"sync"
)

func NoteListener() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		likeNote()
	}()
	go func() {
		defer wg.Done()
		collectNote()
	}()
	log.Printf("note mq listener has started")
	wg.Wait()
}

// LikeNote 点赞&取消点赞帖子
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
			break
		}

		err = json.Unmarshal(message.Value, &msg)
		if err != nil {
			break
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

// 收藏&取消收藏帖子
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
			break
		}

		err = json.Unmarshal(message.Value, &msg)
		if err != nil {
			break
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
