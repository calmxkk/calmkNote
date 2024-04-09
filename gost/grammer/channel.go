package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/util/grand"
)

var queue = make(chan int, 10)

var wg = sync.WaitGroup{}

func worker() {
	for {
		temp := grand.N(1, 100)
		queue <- temp
		fmt.Println("input ", temp)
		time.Sleep(time.Second)
	}
}

func custom() {
	for {
		select {
		case i := <-queue:
			fmt.Println("                recv ", i)
			time.Sleep(time.Second)
		}
	}
}

func test_channel() {
	go worker()
	go worker()
	go worker()
	go custom()

	wg.Add(4)
	wg.Wait()
}
