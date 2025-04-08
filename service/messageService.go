package service

import (
	"log"
	"note_app_server_mq/model/msgModel"
	"note_app_server_mq/repository"
)

func SyncToMongo(firstKey int64, secondKey int64, msg *msgModel.Message) {
	err := repository.SyncMessageToMongo(firstKey, secondKey, msg)
	if err != nil {
		log.Println("failed to sync message:", err)
	}
}
