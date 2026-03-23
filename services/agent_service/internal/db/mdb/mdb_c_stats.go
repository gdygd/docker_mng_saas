package mdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"agent-service/internal/db"
	"agent-service/internal/logger"
)

// InsertContainerStats — 수집 이력 bulk INSERT (시계열, UPSERT 없음)
func (q *MariaDbHandler) InsertContainerStats(ctx context.Context, agentid, hostId int, params []db.ContainerStatsParams) error {
	if len(params) == 0 {
		return nil
	}

	ado := q.GetDB()

	tx, err := ado.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		logger.Log.Error("failed to begin transaction: %v", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
		INSERT INTO container_stats_log (
			id,
			host_id,
			container_id,
			collected_at,
			container_name,
			cpu_percent,
			memory_usage,
			memory_limit,
			memory_percent,
			network_rx,
			network_tx
		) VALUES `

	placeholders := make([]string, 0, len(params))
	args := make([]interface{}, 0, len(params)*10)

	for _, p := range params {
		placeholders = append(placeholders, "(?, ?, ?, NOW(), ?, ?, ?, ?, ?, ?, ?)")
		args = append(args,
			agentid, hostId, p.ID, p.Name,
			p.CPUPercent, p.MemoryUsage, p.MemoryLimit, p.MemoryPercent,
			p.NetworkRx, p.NetworkTx,
		)
	}

	query += strings.Join(placeholders, ", ")

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		logger.Log.Error("failed to insert container stats: %v", err)
		return fmt.Errorf("failed to insert container stats: %w", err)
	}

	if err = tx.Commit(); err != nil {
		logger.Log.Error("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
