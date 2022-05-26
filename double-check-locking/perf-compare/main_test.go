package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

type Conn struct {
	Addr  string
	State int
}

func (c *Conn) Send() {
	// do nothing
}

type Context struct {
	done uint32
	c    *Conn
	mu   sync.Mutex
}

func BenchmarkInfo_UnsafeDCL(b *testing.B) {
	context := &Context{}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ensure_unsafe_dcl(context)
			processConn(context.c)
		}
	})
}

func ensure_unsafe_dcl(context *Context) {
	if context.done == 0 {
		context.mu.Lock()
		defer context.mu.Unlock()
		if context.done == 0 {
			defer func() { context.done = 1 }()
			context.c = newConn()
		}
	}
}

func BenchmarkInfo_DCL(b *testing.B) {
	context := &Context{}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ensure_dcl(context)
			processConn(context.c)
		}
	})
}

func ensure_dcl(context *Context) {
	if atomic.LoadUint32(&context.done) == 0 {
		context.mu.Lock()
		defer context.mu.Unlock()
		if context.done == 0 {
			defer atomic.StoreUint32(&context.done, 1)
			context.c = newConn()
		}
	}
}

func BenchmarkInfo_Mutex(b *testing.B) {
	context := &Context{}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ensure_mutex(context)
			processConn(context.c)
		}
	})
}

func ensure_mutex(context *Context) {
	context.mu.Lock()
	defer context.mu.Unlock()
	if context.done == 0 {
		defer func() { context.done = 1 }()
		context.c = newConn()
	}
}

func newConn() *Conn {
	return &Conn{"127.0.0.1:8080", 1}
}

func processConn(c *Conn) {
	c.Send()
}
