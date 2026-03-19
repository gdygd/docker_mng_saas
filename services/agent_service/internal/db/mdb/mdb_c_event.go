package mdb

import (
	"agent-service/internal/db"
	"agent-service/internal/logger"
	"context"
	"fmt"
)

// InsertContainerEvent — 단건 이벤트 INSERT
func (q *MariaDbHandler) InsertContainerEvent(ctx context.Context, hostId int, param db.ContainerEventParams) error {
	ado := q.GetDB()

	query := `
		INSERT INTO container_event (
			id,
			container_id,
			received_at,
			hostname,
			type,
			action,
			actor_id,
			actor_name,
			event_timestamp,
			attrs
		) VALUES (?, ?, NOW(), ?, ?, ?, ?, ?, ?, ?)`

	_, err := ado.ExecContext(ctx, query,
		hostId, param.ContainerID, param.Hostname,
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
