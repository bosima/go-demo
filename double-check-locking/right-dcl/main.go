package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Conn struct {
	Addr  string
	State int
}

var c *Conn
var mu sync.Mutex
var done uint32
var i = 0

func getInstance() *Conn {
	if atomic.LoadUint32(&done) == 0 {
		mu.Lock()
		defer mu.Unlock()

		if done == 0 {
			defer func() {
				if r := recover(); r == nil {
					defer atomic.StoreUint32(&done, 1)
				}
			}()

			c = mustNewConn()
		}
	}
	return c
}

// throw an panic directly

// func newConn() *Conn {
// 	fmt.Println("newConn")
// 	div := i
// 	i++
// 	k := 1 / div
// 	return &Conn{"127.0.0.1:8080", k}
// }

func mustNewConn() *Conn {
	conn, err := newConn()
	if err != nil {
		panic(err)
	}
	return conn
}

func newConn() (*Conn, error) {
	fmt.Println("newConn")
	div := i
	i++
	if div == 0 {
		return nil, errors.New("the divisor is zero")
	}
	k := 1 / div
	return &Conn{"127.0.0.1:8080", k}, nil
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
