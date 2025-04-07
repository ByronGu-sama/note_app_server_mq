package global

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"sync"
)

var (
	Db               *gorm.DB                   // Db mysql
	ESClient         *elasticsearch.TypedClient // ESClient es客户端
	OssClientPool    *sync.Pool                 // OssClientPool oss客户端连接池
	MongoClient      *mongo.Client              // MongoClient mongoDB客户端
	NoteNormalRdb    *redis.Client              // NoteNormalRdb 缓存笔记相关的附加数据
	CommentNormalRdb *redis.Client              // CommentNormalRdb 缓存评论相关的附加数据
)
