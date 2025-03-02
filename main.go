package main

import (
	"note_app_server_mq/comsumer"
	"note_app_server_mq/config"
)

func main() {
	config.InitAppConfig()
	comsumer.NoteListener()
	//test.TestKafka()
}
