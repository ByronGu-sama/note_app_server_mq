package utils

import "fmt"

func SafeGo(fn func()) {
	go func() {
		if err := recover(); err != nil {
			fmt.Println("recover from panic")
		}
	}()
	fn()
}
