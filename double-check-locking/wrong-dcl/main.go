package main

import (
	"fmt"
	"sync"
	"time"
)

type Conn struct {
	Addr  string
	State int
}

var c *Conn
var mu sync.Mutex

func getInstance() *Conn {
	if c == nil {
		mu.Lock()
		defer mu.Unlock()
		if c == nil {
			c = &Conn{"127.0.0.1:8080", 1}
		}
	}
	return c
}

func doprint() {
	ins := getInstance()
	fmt.Println(ins)
}

func loopPrint() {
	for i := 0; i < 10; i++ {
		go doprint()
	}
}

func main() {
	loopPrint()
	time.Sleep(time.Second * 2)
}
