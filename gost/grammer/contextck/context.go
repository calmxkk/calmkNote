package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	test_deadline()
}

var AllWg = sync.WaitGroup{}

func test_deadline() {
	rootctx := context.Background()
	ctx1, _ := context.WithDeadline(rootctx, time.Now().Add(time.Second*5))
	ctx2, _ := context.WithDeadline(ctx1, time.Now().Add(time.Second*2))

	AllWg.Add(2)
	go work(ctx1, "ctx1")
	go work(ctx2, "ctx2")
	AllWg.Wait()
	return
}

func test_withchannel() {
	rootctx := context.Background()
	ctx1, _ := context.WithCancel(rootctx)
	ctx2, cancel2 := context.WithCancel(ctx1)
	ctx3, cancel3 := context.WithCancel(rootctx)

	AllWg.Add(3)
	go work(ctx1, "ctx1")
	go work(ctx2, "ctx2")
	go work(ctx3, "ctx3")

	time.Sleep(time.Second * 3)
	// cancel1()
	time.Sleep(time.Second * 3)
	cancel2()
	time.Sleep(time.Second * 3)
	cancel3()

	AllWg.Wait()
	return
}

func work(ctx context.Context, name string) {
	ticker := time.NewTicker(time.Second)
	defer AllWg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("grouting %v stop\n", name)
			return
		case t := <-ticker.C:
			fmt.Printf("grouting %v %v\n", name, t)
		}
	}
}
