package dbbatch

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"agent-service/internal/db"
	"agent-service/internal/dto"
	"agent-service/internal/logger"
)

const (
	maxBatchSize  = 200
	flushInterval = 500 * time.Millisecond
	dbTimeout     = 10 * time.Second

	// MariaDB placeholder 한계(65535) 기준 청크 사이즈
	// Info    : 7  cols × 500 = 3,500
	// Inspect : 11 cols × 200 = 2,200
	// Stats   : 10 cols × 300 = 3,000
	defaultChunkSize = 300
	infoChunkSize    = 200
	inspectChunkSize = 200
	statsChunkSize   = 300
	eventChunkSize   = 300
)

// ============================================================================
// Workers
// ============================================================================

func (b *DbBatch) runInfoWorker() {
	defer b.workerWg.Done()

	batch := make([]ContainerInfoItem, 0, maxBatchSize)
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
		defer cancel()

		// 배치 내 전체 항목을 하나의 params 슬라이스로 합산
		allParams := make([]db.ContainerInfoParams, 0, len(batch)*10)

		for _, item := range batch {
			params := toContainerInfoParams(item.AgentId, item.HostId, item.Data)
			allParams = append(allParams, params...)
			// if err := b.dbHnd.CreateContainerInfo(ctx, item.AgentId, item.HostId, params); err != nil {
			// 	logger.Log.Error("DbBatch[info] flush error agent[%d]: %v", item.AgentId, err)
			// }
		}
		batch = batch[:0]

		//----------------------
		// defaultChunkSize 단위로 분할 INSERT — placeholder 한계(65535) 초과 방지
		total := len(allParams)
		chunks := (total + infoChunkSize - 1) / infoChunkSize
		for i := 0; i < total; i += infoChunkSize {
			end := i + infoChunkSize
			if end > total {
				end = total
			}
			logger.Log.Print(1, "insert... : %d", end-i)
			if err := b.dbHnd.CreateContainerInfo(ctx, allParams[i:end]); err != nil {
				logger.Log.Error("DbBatch[info] flush error: %v", err)
			}
			logger.Log.Print(1, "insert ok : %d", end-i)
		}
		logger.Log.Print(1, "info flush: %d rows, %d chunks", total, chunks)

		//----------------------
	}

	for {
		select {
		case item := <-b.infoQ:
			batch = append(batch, item)
			if len(batch) >= maxBatchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-b.ctx.Done():
			// 잔여 항목 드레인
			for draining := true; draining; {
				select {
				case item := <-b.infoQ:
					batch = append(batch, item)
				default:
					draining = false
				}
			}
			flush()
			return
		}
	}
}

func (b *DbBatch) runInspectWorker() {
	defer b.workerWg.Done()

	batch := make([]ContainerInspectItem, 0, maxBatchSize)
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
		defer cancel()

		allParams := make([]db.ContainerInspectParams, 0, len(batch)*10)

		for _, item := range batch {
			params, err := toContainerInspectParams(item.AgentId, item.HostId, item.Data)
			if err != nil {
				logger.Log.Error("DbBatch[inspect] convert error agent[%d]: %v", item.AgentId, err)
				continue
			}
			allParams = append(allParams, params...)

			// if err := b.dbHnd.UpsertContainerInspect(ctx, item.AgentId, item.HostId, params); err != nil {
			// 	logger.Log.Error("DbBatch[inspect] flush error agent[%d]: %v", item.AgentId, err)
			// }
		}
		batch = batch[:0]

		//--------------------
		// defaultChunkSize 단위로 분할 INSERT — placeholder 한계(65535) 초과 방지
		defer cancel()
		total := len(allParams)
		chunks := (total + inspectChunkSize - 1) / inspectChunkSize
		for i := 0; i < total; i += inspectChunkSize {
			end := i + inspectChunkSize
			if end > total {
				end = total
			}
			logger.Log.Print(1, "update inspect... : %d", end-i)
			if err := b.dbHnd.UpsertContainerInspect(ctx, allParams[i:end]); err != nil {
				logger.Log.Error("DbBatch[inspect] flush error: %v", err)
			}
			logger.Log.Print(1, "update inspect ok : %d", end-i)
		}
		logger.Log.Print(1, "inspect flush: %d rows, %d chunks", total, chunks)

		//--------------------
	}

	for {
		select {
		case item := <-b.inspectQ:
			batch = append(batch, item)
			if len(batch) >= maxBatchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-b.ctx.Done():
			for draining := true; draining; {
				select {
				case item := <-b.inspectQ:
					batch = append(batch, item)
				default:
					draining = false
				}
			}
			flush()
			return
		}
	}
}

func (b *DbBatch) runStatsWorker() {
	defer b.workerWg.Done()

	batch := make([]ContainerStatsItem, 0, maxBatchSize)
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		// 배치 내 전체 항목을 하나의 params 슬라이스로 합산
		allParams := make([]db.ContainerStatsParams, 0, len(batch)*10)
		for _, item := range batch {
			allParams = append(allParams, toContainerStatsParams(item.AgentId, item.HostId, item.Data)...)
		}
		batch = batch[:0]

		// statsChunkSize 단위로 분할 INSERT — placeholder 한계(65535) 초과 방지
		ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
		defer cancel()
		total := len(allParams)
		chunks := (total + statsChunkSize - 1) / statsChunkSize
		for i := 0; i < total; i += statsChunkSize {
			end := i + statsChunkSize
			if end > total {
				end = total
			}
			logger.Log.Print(1, "insert... : %d", end-i)
			if err := b.dbHnd.InsertContainerStats(ctx, allParams[i:end]); err != nil {
				logger.Log.Error("DbBatch[stats] flush error: %v", err)
			}
			logger.Log.Print(1, "insert ok : %d", end-i)
		}
		logger.Log.Print(1, "stats flush: %d rows, %d chunks", total, chunks)
	}

	for {
		select {
		case item := <-b.statsQ:
			batch = append(batch, item)
			if len(batch) >= maxBatchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-b.ctx.Done():
			for draining := true; draining; {
				select {
				case item := <-b.statsQ:
					batch = append(batch, item)
				default:
					draining = false
				}
			}
			flush()
			return
		}
	}
}

func (b *DbBatch) runEventWorker() {
	defer b.workerWg.Done()

	// 이벤트는 단건 INSERT (sequence 사용) — 배치 병합 없이 순차 처리
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case item := <-b.eventQ:
			ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
			param, err := toContainerEventParam(item.AgentId, item.HostId, item.Data)
			if err != nil {
				logger.Log.Error("DbBatch[event] convert error agent[%d]: %v", item.AgentId, err)
				cancel()
				continue
			}
			if err := b.dbHnd.InsertContainerEvent(ctx, param); err != nil {
				logger.Log.Error("DbBatch[event] flush error agent[%d]: %v", item.AgentId, err)
			}
			cancel()
		case <-ticker.C:
			// 이벤트 워커는 ticker 불필요하나 select 균형을 위해 유지
		case <-b.ctx.Done():
			for draining := true; draining; {
				select {
				case item := <-b.eventQ:
					ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
					param, err := toContainerEventParam(item.AgentId, item.HostId, item.Data)
					if err != nil {
						logger.Log.Error("DbBatch[event] drain convert error: %v", err)
						cancel()
						continue
					}
					if err := b.dbHnd.InsertContainerEvent(ctx, param); err != nil {
						logger.Log.Error("DbBatch[event] drain flush error: %v", err)
					}
					cancel()
				default:
					draining = false
				}
			}
			return
		}
	}
}

func (b *DbBatch) runEventWorker2() {
	defer b.workerWg.Done()

	batch := make([]ContainerEventItem, 0, maxBatchSize)
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		// 배치 내 전체 항목을 하나의 params 슬라이스로 합산
		allParams := make([]db.ContainerEventParams, 0, len(batch)*10)
		for _, item := range batch {
			param, _ := toContainerEventParam(item.AgentId, item.HostId, item.Data)
			allParams = append(allParams, param)
		}
		batch = batch[:0]

		// statsChunkSize 단위로 분할 INSERT — placeholder 한계(65535) 초과 방지
		ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
		defer cancel()
		total := len(allParams)
		chunks := (total + eventChunkSize - 1) / eventChunkSize
		for i := 0; i < total; i += eventChunkSize {
			end := i + eventChunkSize
			if end > total {
				end = total
			}
			logger.Log.Print(1, "insert... : %d", end-i)
			if err := b.dbHnd.InsertContainerEvent2(ctx, allParams[i:end]); err != nil {
				logger.Log.Error("DbBatch[stats] flush error: %v", err)
			}
			logger.Log.Print(1, "insert ok : %d", end-i)
		}
		logger.Log.Print(1, "event flush: %d rows, %d chunks", total, chunks)
	}

	for {
		select {
		case item := <-b.eventQ:
			batch = append(batch, item)
			if len(batch) >= maxBatchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-b.ctx.Done():
			for draining := true; draining; {
				select {
				case item := <-b.eventQ:
					batch = append(batch, item)
				default:
					draining = false
				}
			}
			flush()
			return
		}
	}
}

// ============================================================================
// dto → db.Params 변환 (rpcservice.go 와 동일 로직, dbbatch 패키지 전용)
// ============================================================================

func toContainerInfoParams(agentid, hostid int, req dto.ContainerListData) []db.ContainerInfoParams {
	params := make([]db.ContainerInfoParams, 0, len(req.Containers))
	for _, r := range req.Containers {
		params = append(params, db.ContainerInfoParams{
			AgentId: agentid,
			HostId:  hostid,
			ID:      r.ID,
			Name:    r.Name,
			Image:   r.Image,
			State:   r.State,
			Status:  r.Status,
		})
	}
	return params
}

func toContainerInspectParams(agentid, hostid int, req dto.ContainerInspectData) ([]db.ContainerInspectParams, error) {
	params := make([]db.ContainerInspectParams, 0, len(req.Inspects))
	for _, r := range req.Inspects {
		stateJSON, err := marshalJSON(r.State)
		if err != nil {
			return nil, fmt.Errorf("state marshal: %w", err)
		}
		configJSON, err := marshalJSON(r.Config)
		if err != nil {
			return nil, fmt.Errorf("config marshal: %w", err)
		}
		networkJSON, err := marshalJSON(r.Network)
		if err != nil {
			return nil, fmt.Errorf("network marshal: %w", err)
		}
		mountJSON, err := marshalJSON(r.Mounts)
		if err != nil {
			return nil, fmt.Errorf("mount marshal: %w", err)
		}
		params = append(params, db.ContainerInspectParams{
			AgentId:      agentid,
			HostId:       hostid,
			ID:           r.ID,
			Name:         r.Name,
			Image:        r.Image,
			Platform:     r.Platform,
			RestartCount: r.RestartCount,
			StateInfo:    stateJSON,
			ConfigInfo:   configJSON,
			NetworkInfo:  networkJSON,
			MountInfo:    mountJSON,
		})
	}
	return params, nil
}

func toContainerStatsParams(agentid, hostid int, req dto.ContainerStatsData) []db.ContainerStatsParams {
	params := make([]db.ContainerStatsParams, 0, len(req.Stats))
	for _, r := range req.Stats {
		params = append(params, db.ContainerStatsParams{
			AgentId:       agentid,
			HostId:        hostid,
			ID:            r.ID,
			Name:          r.Name,
			CPUPercent:    r.CPUPercent,
			MemoryUsage:   r.MemoryUsage,
			MemoryLimit:   r.MemoryLimit,
			MemoryPercent: r.MemoryPercent,
			NetworkRx:     r.NetworkRx,
			NetworkTx:     r.NetworkTx,
		})
	}
	return params
}

func toContainerEventParam(agentid, hostid int, req dto.ContainerEvent) (db.ContainerEventParams, error) {
	attrsJSON, err := marshalJSON(req.Attrs)
	if err != nil {
		return db.ContainerEventParams{}, fmt.Errorf("attrs marshal: %w", err)
	}
	return db.ContainerEventParams{
		AgentId:     agentid,
		HostId:      hostid,
		ContainerID: req.ActorID,
		Hostname:    req.Host,
		Type:        req.Type,
		Action:      req.Action,
		ActorID:     req.ActorID,
		ActorName:   req.ActorName,
		EventTime:   req.Timestamp,
		Attrs:       attrsJSON,
	}, nil
}

func marshalJSON(v any) (string, error) {
	if v == nil {
		return "null", nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
