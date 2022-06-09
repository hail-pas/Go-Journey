package main

import (
	"fmt"
	"runtime"
)

func say(s string) {
	runtime.Gosched()
	for i := 0; i < 5; i++ {
		fmt.Println(s)
	}
}

func main() {
	go say("world")
	say("hello")
}
