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
)

func CommentListener() {
	log.Printf("comment mq listener has started")
	likeComment()
	delComment()
}

func delComment() {
	ctx := context.Background()
	conn, err := kafka.DialLeader(ctx, config.AC.Kafka.Network, config.AC.Kafka.Host+":"+config.AC.Kafka.Port, config.AC.Kafka.NoteComments.Topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	for {
		msg := mqMessageModel.DelNoteComment{}
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
		cid := msg.Cid
		if msg.Action == action.DelNoteComment {
			err = repository.DeleteComment(uid, cid)
			if err != nil {
				log.Println("failed to like note:", err)
			}
		}
	}
}

func likeComment() {
	ctx := context.Background()
	conn, err := kafka.DialLeader(ctx, config.AC.Kafka.Network, config.AC.Kafka.Host+":"+config.AC.Kafka.Port, config.AC.Kafka.NoteComments.Topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	for {
		msg := mqMessageModel.LikeNoteComment{}
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
