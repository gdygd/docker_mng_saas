package gapi

import (
	"context"
	"sync/atomic"
	"time"

	"agent-service/internal/logger"
)

// statsCounter — ContainerStats RPC 초당 수신 건수 카운터
type statsCounter struct {
	count int64 // atomic
}

func (c *statsCounter) inc() {
	atomic.AddInt64(&c.count, 1)
}

func (c *statsCounter) swapAndGet() int64 {
	return atomic.SwapInt64(&c.count, 0)
}

// runStatsRateLogger — 1초마다 수신 건수를 출력하는 백그라운드 고루틴
func (server *Server) runStatsRateLogger(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			n := server.statsCounter.swapAndGet()
			logger.Log.Print(1, "ContainerStats RPS  %s : %d", t.Format("15:04:05"), n)
		case <-ctx.Done():
			return
		}
	}
}
