package appModel

import "time"

type AppConfig struct {
	App struct {
		Name string
		Host string
		Port string
	}
	Mysql struct {
		Dsn             string
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime string
	}
	Redis struct {
		NoteNormalRdb    int
		CommentNormalRdb int
		Host             string
		Port             string
		Password         string
		Timeout          time.Duration
		Pool             struct {
			MaxActive int
			MaxIdle   int
			MinIdle   int
			MaxWait   time.Duration
		}
	}
	Es struct {
		Host string
		Port string
	}
	Mongo struct {
		Host     string
		Port     string
		Username string
		Password string
	}
	Kafka struct {
		Network   string
		Host      string
		Port      string
		NoteLikes struct {
			Topic      string
			Partitions int
		}
		NoteCollects struct {
			Topic      string
			Partitions int
		}
		NoteComments struct {
			Topic      string
			Partitions int
		}
		SyncNotes struct {
			Topic      string
			Partitions int
		}
		DelNotes struct {
			Topic      string
			Partitions int
		}
		SyncMessages struct {
			Topic      string
			Partitions int
		}
	}
	Oss struct {
		AvatarBucket   string
		NotePicsBucket string
		StyleBucket    string
		EndPoint       string
		Region         string
	}
}
