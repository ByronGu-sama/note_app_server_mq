package config

import (
	"github.com/redis/go-redis/v9"
	"note_app_server_mq/global"
	"time"
)

func InitRedisConfig() {
	noteNormalRdb := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.NoteNormalRdb,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	commentNormalRdb := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.CommentNormalRdb,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	global.NoteNormalRdb = noteNormalRdb
	global.CommentNormalRdb = commentNormalRdb
}
