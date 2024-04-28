package main

import (
	"context"
	"fmt"
	"sync"
)

var serverwg sync.WaitGroup

func worker(ctx context.Context) {
	serverwg.Add(1)
	defer serverwg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("recv context done, grouting exit")
		default:
			fmt.Println("working")
		}
	}
}
