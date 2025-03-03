package global

import (
	"github.com/elastic/go-elasticsearch/v8"
	"gorm.io/gorm"
	"sync"
)

var (
	Db            *gorm.DB                   // Db mysql
	ESClient      *elasticsearch.TypedClient // ESClient es客户端
	OssClientPool *sync.Pool                 // OssClientPool oss客户端连接池
)
