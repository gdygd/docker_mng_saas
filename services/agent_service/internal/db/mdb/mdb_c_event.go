package mdb

import (
	"context"
	"fmt"

	"agent-service/internal/db"
	"agent-service/internal/logger"
)

// InsertContainerEvent — 단건 이벤트 INSERT
func (q *MariaDbHandler) InsertContainerEvent(ctx context.Context, agentid, hostId int, param db.ContainerEventParams) error {
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
		agentid, hostId, param.ContainerID, param.Hostname,
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
