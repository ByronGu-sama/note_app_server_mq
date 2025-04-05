package writer

import (
	"time"
)

// 定期向数据库内写入缓存中的点赞数
func writeThumbsUpToDB(gap time.Duration) {
	// 定期将缓存数据写入DB
	for {
		select {
		case <-time.Tick(gap):

		}
	}
}
