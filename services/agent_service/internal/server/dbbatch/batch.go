package dbbatch

import (
	"context"
	"sync"

	"agent-service/internal/container"
	"agent-service/internal/db"
	"agent-service/internal/dto"
	"agent-service/internal/logger"
)

const (
	queueSize        = 30000
	statsWorkerCount = 5 * 3
)

// DbBatch — 비동기 배치 DB 쓰기 처리기
// gRPC 핸들러로부터 dto 데이터를 받아 채널에 적재하고,
// 백그라운드 워커가 context.Background()로 DB에 순차 기록한다.
type DbBatch struct {
	ctx      context.Context
	cancel   context.CancelFunc
	workerWg sync.WaitGroup
	dbHnd    db.DbHandler

	infoQ    chan ContainerInfoItem
	inspectQ chan ContainerInspectItem
	statsQ   chan ContainerStatsItem
	eventQ   chan ContainerEventItem
}

func NewDbBatch(ct *container.Container) *DbBatch {
	ctx, cancel := context.WithCancel(context.Background())
	return &DbBatch{
		ctx:      ctx,
		cancel:   cancel,
		dbHnd:    ct.DbHnd,
		infoQ:    make(chan ContainerInfoItem, queueSize),
		inspectQ: make(chan ContainerInspectItem, queueSize),
		statsQ:   make(chan ContainerStatsItem, queueSize),
		eventQ:   make(chan ContainerEventItem, queueSize),
	}
}

// Start — 워커 고루틴 실행 (non-blocking)
// Stats는 순수 INSERT(시계열)이므로 동일 채널을 복수 워커가 경쟁적으로 소비해도 안전
func (b *DbBatch) Start() {
	logger.Log.Print(2, "DbBatch workers starting... (stats workers: %d)", statsWorkerCount)
	b.workerWg.Add(2 + statsWorkerCount + 1) // info + inspect + stats×N + event
	go b.runInfoWorker()
	go b.runInspectWorker()
	for i := 0; i < statsWorkerCount; i++ {
		go b.runStatsWorker()
	}
	go b.runEventWorker()
}

// Shutdown — ctx 취소 후 워커가 잔여 큐를 소진할 때까지 대기
func (b *DbBatch) Shutdown() {
	logger.Log.Print(2, "DbBatch shutdown: draining queues...")
	b.cancel()
	b.workerWg.Wait()
	logger.Log.Print(2, "DbBatch shutdown: complete")
}

// ============================================================================
// Push methods — non-blocking. 큐가 가득 찬 경우 드롭 후 로깅
// ============================================================================

func (b *DbBatch) PushContainerInfo(agentId, hostId int, data dto.ContainerListData) {
	select {
	case b.infoQ <- ContainerInfoItem{AgentId: agentId, HostId: hostId, Data: data}:
	default:
		logger.Log.Error("DbBatch: info queue full, dropping agent[%d]", agentId)
	}
}

func (b *DbBatch) PushContainerInspect(agentId, hostId int, data dto.ContainerInspectData) {
	select {
	case b.inspectQ <- ContainerInspectItem{AgentId: agentId, HostId: hostId, Data: data}:
	default:
		logger.Log.Error("DbBatch: inspect queue full, dropping agent[%d]", agentId)
	}
}

func (b *DbBatch) PushContainerStats(agentId, hostId int, data dto.ContainerStatsData) {
	select {
	case b.statsQ <- ContainerStatsItem{AgentId: agentId, HostId: hostId, Data: data}:
		// if len(b.statsQ) > 800 {
		// 	logger.Log.Print(2, "statsQ : %d", len(b.statsQ))
		// }
	default:
		logger.Log.Error("DbBatch: stats queue full, dropping agent[%d]", agentId)
	}
}

// PushContainerEvent — 이벤트는 드롭 없이 블로킹 시도 (중요 데이터)
// ctx 취소 시에만 포기한다.
func (b *DbBatch) PushContainerEvent(agentId, hostId int, data dto.ContainerEvent) {
	select {
	case b.eventQ <- ContainerEventItem{AgentId: agentId, HostId: hostId, Data: data}:
	case <-b.ctx.Done():
		logger.Log.Error("DbBatch: shutdown in progress, dropping event agent[%d]", agentId)
	}
}
