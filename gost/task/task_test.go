package task

import (
	"context"
	"testing"
)

func TestGetMemory(t *testing.T) {
	tempTask := NewDockerTask()
	tempTask.GetSystemInfo(context.Background())
}
