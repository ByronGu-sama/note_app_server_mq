package appModel

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
	}
}
