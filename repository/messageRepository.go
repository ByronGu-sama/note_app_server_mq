package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"note_app_server_mq/global"
	"note_app_server_mq/model/msgModel"
)

// SyncMessageToMongo 同步消息
func SyncMessageToMongo(ctx context.Context, uid1, uid2 int64, message *msgModel.Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		client := global.MongoClient
		msg := bson.D{
			{Key: "uid1", Value: uid1},
			{Key: "uid2", Value: uid2},
			{Key: "fromId", Value: message.FromId},
			{Key: "fromName", Value: message.FromName},
			{Key: "toId", Value: message.ToId},
			{Key: "toName", Value: message.ToName},
			{Key: "fromAvatar", Value: message.FromAvatar},
			{Key: "type", Value: message.Type},
			{Key: "content", Value: message.Content},
			{Key: "mediaType", Value: message.MediaType},
			{Key: "url", Value: message.Url},
			{Key: "pubTime", Value: message.PubTime},
			{Key: "groupId", Value: message.GroupId},
			{Key: "read", Value: message.Read},
		}
		messageCollection := client.Database("pending_message").Collection("msgs")
		_, err := messageCollection.InsertOne(context.TODO(), msg)
		if err != nil {
			return err
		}
		return nil
	}
}
