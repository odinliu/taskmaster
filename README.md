# taskmaster
[![GoDoc](https://godoc.org/github.com/odinliu/taskmaster?status.png)](https://godoc.org/github.com/odinliu/taskmaster)

A simple supervised goroutine pool based on worker-thread model.

# Go version
test on go v1.5.1

# Installation
go get github.com/odinliu/taskmaster

# Usage
See godoc [here](https://godoc.org/github.com/odinliu/taskmaster)

# Example
```
package main

import (
	"fmt"
	"time"

	"github.com/odinliu/taskmaster"
)

type MyLogger struct{}

func (MyLogger) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func panicable() {
	select {
	case <-time.After(5 * time.Second):
		panic("paniced!")
	}
}

func main() {
	master := taskmaster.NewSupervisor(func() {
		panicable()
	}, taskmaster.SuperOption{
		NeedRestart:    true,
		RestartDelay:   5 * time.Second,
		MaxFailureTime: 1 * 100,
		MaxWorkerNum:   1,
		Logger:         &MyLogger{},
	})
	master.Start()
	// infinite loop, quit with ctrl+c
	for true {
	}
}
```
