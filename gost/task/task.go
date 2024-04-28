package task

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"gost/procutil"
	"time"
)

type DockerTask struct {
	timeDuration time.Duration
}

func NewDockerTask() *DockerTask {
	return &DockerTask{
		timeDuration: time.Second * 5,
	}
}

func (d *DockerTask) Run(ctx context.Context) {
	go d.GetMemory(ctx)
	go d.GetSystemInfo(ctx)
}

func (d *DockerTask) GetMemory(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			g.Log().Info(ctx, "grouting getMemory exit")
			return
		default:
			res, err := procutil.GetMemoryInfo(ctx)
			if err != nil {
				g.Log().Errorf(ctx, "%v", err.Error())
				return
			}

			fmt.Printf("%s memory: total %v, used %v\n", time.Now().Format("2006-01-02 15:04:05"), procutil.FmtByte(int64(res.Total)), procutil.FmtByte(int64(res.Used)))
			time.Sleep(d.timeDuration)
		}
	}
}

func (d *DockerTask) GetSystemInfo(ctx context.Context) {
	for {
		select {
		default:
			res, err := procutil.GetSystemInfo(ctx)
			if err != nil {
				return
			}

			fmt.Printf("%s get system Info %v\n", time.Now().Format("2006-01-02 15:04:05"), res)
			time.Sleep(d.timeDuration)
		case <-ctx.Done():
			g.Log().Info(ctx, "grouting getMemory exit")
			return
		}
	}
}
