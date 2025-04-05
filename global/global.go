package global

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"sync"
)

var (
	Db                          *gorm.DB                   // Db mysql
	ESClient                    *elasticsearch.TypedClient // ESClient es客户端
	OssClientPool               *sync.Pool                 // OssClientPool oss客户端连接池
	MongoClient                 *mongo.Client              // MongoClient mongoDB客户端
	ThumbsUpRdbClient           *redis.Client              // ThumbsUpRdbClient ThumbsUpRdbClient缓存点赞数据
	UserLikedNotesRdbClient     *redis.Client              // UserLikedNotesRdbClient ThumbsUpRdbClient缓存点赞数据
	CollectedCntRdbClient       *redis.Client              // CollectedCntRdbClient ThumbsUpRdbClient缓存点赞数据
	UserCollectedNotesRdbClient *redis.Client              // UserCollectedNotesRdbClient ThumbsUpRdbClient缓存点赞数据
)
