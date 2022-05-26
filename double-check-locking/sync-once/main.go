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
var once sync.Once

func setInstance() {
	fmt.Println("setup")
	c = &Conn{"127.0.0.1:8080", 1}
}

func getInstance() *Conn {
	once.Do(setInstance)
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
