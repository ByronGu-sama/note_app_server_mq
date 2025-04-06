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

func CommentListener() {
	fmt.Println("comment MQ listener has started")
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		likeComment()
	}()
	go func() {
		defer wg.Done()
		delComment()
	}()
	wg.Wait()
}

// 删除评论
func delComment() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.AC.Kafka.Host + ":" + config.AC.Kafka.Port},
		GroupID:  "comment_del_group",
		Topic:    config.AC.Kafka.NoteComments.Topic,
		MaxBytes: 10e3,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		message, err1 := reader.ReadMessage(ctx)

		if err1 != nil {
			log.Fatal("failed to read message:", err1)
		}

		msg := mqMessageModel.DelNoteComment{}
		err := json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Fatal("failed to unmarshal message:", err)
		}

		log.Println(fmt.Sprintf("Fetched New Msg:%v", msg))

		uid := msg.Uid
		cid := msg.Cid
		if msg.Action == action.DelNoteComment {
			err = repository.DeleteComment(uid, cid)
			if err != nil {
				log.Println("failed to like note:", err)
			}
		}
	}
}

// 点赞评论
func likeComment() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.AC.Kafka.Host + ":" + config.AC.Kafka.Port},
		GroupID:  "comment_thumbsUp_group",
		Topic:    config.AC.Kafka.NoteComments.Topic,
		MaxBytes: 10e3,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		message, err1 := reader.ReadMessage(ctx)
		if err1 != nil {
			log.Fatal("failed to read message:", err1)
		}

		msg := mqMessageModel.LikeNoteComment{}
		err := json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Fatal("failed to unmarshal message:", err)
		}

		log.Println(fmt.Sprintf("Fetched New Msg:%v", msg))

		uid := msg.Uid
		cid := msg.Cid
		if msg.Action == action.LikeComment {
			err = repository.LikeComment(uid, cid)
			if err != nil {
				log.Println("failed to like note:", err)
			}
		} else if msg.Action == action.DislikeComment {
			err = repository.DislikeComment(uid, cid)
			if err != nil {
				log.Println("failed to dislike note:", err)
			}
		}
	}
}
