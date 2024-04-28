package procutil

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shirou/gopsutil/v3/host"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

type DiskInfo struct {
	Total   uint64
	Percent float64
	Used    uint64
}

type SystemInfo struct {
	Uptime          string
	Hostname        string
	Os              string
	Platform        string
	PlatformVersion string
	KernelVersion   string
	KernelArch      string
}

type MemoryInfo struct {
	Total   uint64
	Percent float64
	Used    uint64
}

func GetMemoryInfo(ctx context.Context) (*MemoryInfo, error) {
	v, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "get memory error, %v\n", err.Error())
		return nil, err
	}

	return &MemoryInfo{
		Total:   v.Total,
		Percent: v.UsedPercent,
		Used:    v.Used,
	}, nil
}

func GetSystemInfo(ctx context.Context) (*SystemInfo, error) {
	sysInfo, err := host.InfoWithContext(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "get system info error, %v\n", err.Error())
		return nil, err
	}

	return &SystemInfo{
		Uptime:          time.Unix(int64(sysInfo.BootTime), 0).Local().Format("2006-01-02 15:04:05"),
		Hostname:        sysInfo.Hostname,
		Os:              sysInfo.OS,
		Platform:        sysInfo.Platform,
		PlatformVersion: sysInfo.PlatformVersion,
		KernelVersion:   sysInfo.KernelVersion,
		KernelArch:      sysInfo.KernelArch,
	}, nil
}

// 字节单位转换
func FmtByte(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%.2fB", float64(size)/float64(1))
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2fKB", float64(size)/float64(1024))
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2fMB", float64(size)/float64(1024*1024))
	} else if size < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2fGB", float64(size)/float64(1024*1024*1024))
	} else if size < 1024*1024*1024*1024*1024 {
		return fmt.Sprintf("%.2fTB", float64(size)/float64(1024*1024*1024*1024))
	} else {
		return fmt.Sprintf("%.2fEB", float64(size)/float64(1024*1024*1024*1024*2014))
	}
}
