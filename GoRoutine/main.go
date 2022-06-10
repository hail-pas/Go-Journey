package main

import (
	"fmt"
	"sync"
	"time"
)

var x int64

var waitGroup sync.WaitGroup

var lock sync.Mutex

//var rwLock sync.RWMutex

func write() {
	lock.Lock()
	x = x + 1
	time.Sleep(2 * time.Second)
	lock.Unlock()
	waitGroup.Done()
}

func read() {
	lock.Lock()
	time.Sleep(time.Second)
	lock.Unlock()
	waitGroup.Done()
}

//func main() {
// waitGroup、lock
//start := time.Now()
//waitGroup.Add(200)
//for i := 0; i < 100; i++ {
//	go write()
//}
//
//for i := 0; i < 100; i++ {
//	go read()
//}
//
//waitGroup.Wait()
//
//fmt.Println(time.Now().Sub(start))
//}

var icons map[string]string

var initializeOnce = sync.Once{}

func initializeIcons() {
	icons = map[string]string{
		"left":  "OK1",
		"up":    "OK2",
		"right": "OK3",
		"down":  "OK4",
	}
}

func Icon(name string) string {
	initializeOnce.Do(initializeIcons)
	return icons[name]
}

var syncMap = sync.Map{}

func main() {
	fmt.Println(Icon("left"))
	syncMap.Store("a", "b")
	fmt.Println(syncMap.Load("a"))
}
