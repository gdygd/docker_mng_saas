package dbbatch

import (
	"context"
	"time"

	"agent-service/internal/logger"
)

func (b *DbBatch) runQueueCountLogger(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			infoQcnt := len(b.infoQ)
			ispQcnt := len(b.inspectQ)
			statsQcnt := len(b.statsQ)
			evtQcnt := len(b.eventQ)
			logger.Log.Print(1, "#############>> [%s] info : %d, inspect : %d, stats : %d, event : %d",
				t.Format("15:04:05"), infoQcnt, ispQcnt, statsQcnt, evtQcnt)

		case <-ctx.Done():
			return
		}
	}
}
