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

func MessageListener() {
	fmt.Println("message MQ listener has started")
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		syncMessage()
	}()
	wg.Wait()
}

// 点赞评论
func syncMessage() {
	ctx := context.Background()
	conn, err := kafka.DialLeader(ctx, config.AC.Kafka.Network, config.AC.Kafka.Host+":"+config.AC.Kafka.Port, config.AC.Kafka.SyncMessages.Topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	for {
		msg := mqMessageModel.SyncMessageMsg{}
		message, err1 := conn.ReadMessage(10e3)
		if err1 != nil {
			break
		}

		err = json.Unmarshal(message.Value, &msg)
		if err != nil {
			break
		}

		fmt.Println(msg)

		mid := msg.Mid
		if msg.Action == action.SyncMessage {
			err = repository.SyncMessageToMongo(mid, msg.Message)
			if err != nil {
				log.Println("failed to sync message:", err)
			}
		}
	}
}
