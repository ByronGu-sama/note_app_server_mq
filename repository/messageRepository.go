package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"note_app_server_mq/global"
	"note_app_server_mq/model/msgModel"
)

// SyncMessageToMongo 同步消息
func SyncMessageToMongo(mid string, message *msgModel.Message) error {
	client := global.MongoClient
	msg := bson.D{
		{Key: "_id", Value: mid},
		{Key: "fromId", Value: message.FromId},
		{Key: "toId", Value: message.ToId},
		{Key: "fromAvatar", Value: message.FromAvatar},
		{Key: "toId", Value: message.ToId},
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
