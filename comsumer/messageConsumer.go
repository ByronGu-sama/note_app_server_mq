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
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.AC.Kafka.Host + ":" + config.AC.Kafka.Port},
		GroupID:  "message_sync_group",
		Topic:    config.AC.Kafka.SyncMessages.Topic,
		MaxBytes: 10e3,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		message, err1 := reader.ReadMessage(ctx)

		if err1 != nil {
			log.Fatal("failed to read message:", err1)
		}

		msg := mqMessageModel.SyncMessageMsg{}
		err := json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Fatal("failed to unmarshal message:", err)
		}

		log.Println(fmt.Sprintf("Fetched New Msg:%v", msg))

		if msg.Action == action.SyncMessage {
			err = repository.SyncMessageToMongo(msg.FirstKey, msg.SecondKey, msg.Message)
			if err != nil {
				log.Println("failed to sync message:", err)
			}
		}
	}
}
