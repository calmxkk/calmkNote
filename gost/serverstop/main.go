package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gogf/gf/v2/os/gproc"

	"github.com/gogf/gf/v2/frame/g"
)

var serverWg sync.WaitGroup

func worker(ctx context.Context, i int) {
	serverWg.Add(1)
	defer serverWg.Done()

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("receive context done, grouting %v exit\n", i)
			return
		case <-ticker.C:
			fmt.Println("working")
		}
	}
}

func SignalListen(ctx context.Context, handler ...gproc.SigHandler) {

}

func main() {
	var ctx context.Context
	ctx = context.Background()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	server := g.Server()

	subCtx, cancel := context.WithCancel(ctx)

	for i := 0; i < 3; i++ {
		go worker(subCtx, i)
	}

	go func() {
		sig := <-signalChan
		fmt.Printf("receive signal from system, %v", sig)
		cancel()
		time.Sleep(time.Second * 5)
		_ = server.Shutdown()
	}()

	server.Run()

}
