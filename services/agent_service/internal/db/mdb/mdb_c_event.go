package mdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"agent-service/internal/db"
	"agent-service/internal/logger"
)

// InsertContainerEvent — 단건 이벤트 INSERT
func (q *MariaDbHandler) InsertContainerEvent(ctx context.Context, param db.ContainerEventParams) error {
	ado := q.GetDB()

	query := `
		INSERT INTO container_event_log (
			id,
			host_id,
			container_id,
			received_at,
			seq,
			hostname,
			type,
			action,
			actor_id,
			actor_name,
			event_timestamp,
			attrs
		) VALUES (?, ?, ?, NOW(), NEXT VALUE FOR sq_eventlog, ?, ?, ?, ?, ?, ?, ?)`

	_, err := ado.ExecContext(ctx, query,
		param.AgentId, param.HostId, param.ContainerID, param.Hostname,
		param.Type, param.Action,
		param.ActorID, param.ActorName,
		param.EventTime, param.Attrs,
	)
	if err != nil {
		logger.Log.Error("failed to insert container event: %v", err)
		return fmt.Errorf("failed to insert container event: %w", err)
	}

	return nil
}

func (q *MariaDbHandler) InsertContainerEvent2(ctx context.Context, params []db.ContainerEventParams) error {
	if len(params) == 0 {
		return nil
	}

	ado := q.GetDB()

	tx, err := ado.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		logger.Log.Error("[InsertContainerEvent2]failed to begin transaction: %v", err)
		return fmt.Errorf("[InsertContainerEvent2]failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
		INSERT INTO container_event_log (
			id,
			host_id,
			container_id,
			received_at,
			seq,
			hostname,
			type,
			action,
			actor_id,
			actor_name,
			event_timestamp,
			attrs
		) VALUES `

	placeholders := make([]string, 0, len(params))
	args := make([]interface{}, 0, len(params)*10)

	for _, p := range params {
		placeholders = append(placeholders, "(?, ?, ?, NOW(), NEXT VALUE FOR sq_eventlog, ?, ?, ?, ?, ?, ?, ?)")
		args = append(args,
			p.AgentId, p.HostId, p.ContainerID, p.Hostname,
			p.Type, p.Action,
			p.ActorID, p.ActorName,
			p.EventTime, p.Attrs,
		)
	}

	query += strings.Join(placeholders, ", ")

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		logger.Log.Error("[InsertContainerEvent2]failed to insert container stats: %v", err)
		return fmt.Errorf("[InsertContainerEvent2]failed to insert container stats: %w", err)
	}

	if err = tx.Commit(); err != nil {
		logger.Log.Error("[InsertContainerEvent2]failed to commit transaction: %v", err)
		return fmt.Errorf("[InsertContainerEvent2]failed to commit transaction: %w", err)
	}

	return nil
}
