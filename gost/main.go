package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"gost/task"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var ctx context.Context
	ctx = context.Background()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 创建一个通道用于通知清理工作完成
	cleanupDone := make(chan bool, 1)

	server := g.Server()

	subCtx, cancel := context.WithCancel(ctx)

	tempTask := task.NewDockerTask()
	tempTask.Run(subCtx)

	go func() {
		sig := <-signalChan
		fmt.Printf("receive signal from system, %v\n", sig)
		cancel()
		time.Sleep(time.Second * 5)
		_ = server.Shutdown()
		cleanupDone <- true
	}()

	server.Run()
	<-cleanupDone
	fmt.Println("清理工作完成，程序退出")
}
