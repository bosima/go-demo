package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type User struct {
	Id   int64
	mu   sync.Mutex
	done uint32
	c    *Conn
}

type Conn struct {
	Addr  string
	State int
}

var seq int64

func getInstance(user *User) *Conn {
	if atomic.LoadUint32(&user.done) == 0 {
		user.mu.Lock()
		defer user.mu.Unlock()

		if user.done == 0 {
			defer func() {
				if r := recover(); r == nil {
					defer atomic.StoreUint32(&user.done, 1)
				}
			}()

			user.c = newConn()
		}
	}
	return user.c
}

func newConn() *Conn {
	curSeq := atomic.AddInt64(&seq, 1)
	fmt.Println("newConn")
	return &Conn{fmt.Sprintf("127.0.0.1:%d", curSeq), 1}
}

func doprint(user *User) {
	ins := getInstance(user)
	fmt.Println(user.Id, ins)
}

func loopPrint() {
	users := [10]*User{}
	for i := 0; i < 10; i++ {
		user := &User{Id: int64(i)}
		users[i] = user
	}

	for _, user := range users {
		go doprint(user)
	}

	for _, user := range users {
		go doprint(user)
	}

	for _, user := range users {
		go doprint(user)
	}
}

func main() {
	loopPrint()
	time.Sleep(time.Second * 2)
}
