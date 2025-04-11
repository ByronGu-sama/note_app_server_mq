package service

import (
	"context"
	"log"
	"note_app_server_mq/model/msgModel"
	"note_app_server_mq/repository"
)

func SyncToMongo(ctx context.Context, firstKey int64, secondKey int64, msg *msgModel.Message) {
	select {
	case <-ctx.Done():
		return
	default:
		err := repository.SyncMessageToMongo(ctx, firstKey, secondKey, msg)
		if err != nil {
			log.Println("failed to sync message:", err)
		}
	}
}
