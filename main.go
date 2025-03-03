package main

import (
	"note_app_server_mq/comsumer"
	"note_app_server_mq/config"
	"sync"
)

func main() {
	config.InitAppConfig()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		comsumer.NoteListener()
	}()
	go func() {
		defer wg.Done()
		comsumer.CommentListener()
	}()
	wg.Wait()
	//test.TestKafka()
}
