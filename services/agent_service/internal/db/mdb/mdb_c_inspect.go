package mdb

import (
	"agent-service/internal/db"
	"agent-service/internal/logger"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// UpsertContainerInspect — (id, container_id) UNIQUE KEY 기준 UPSERT
func (q *MariaDbHandler) UpsertContainerInspect(ctx context.Context, hostId int, params []db.ContainerInspectParams) error {
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
		INSERT INTO container_inspect (
			id,
			container_id,
			container_name,
			image,
			platform,
			restart_count,
			state_info,
			config_info,
			network_info,
			mount_info,
			changed_at
		) VALUES `

	placeholders := make([]string, 0, len(params))
	args := make([]interface{}, 0, len(params)*10)

	for _, p := range params {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())")
		args = append(args,
			hostId, p.ID, p.Name, p.Image, p.Platform, p.RestartCount,
			p.StateInfo, p.ConfigInfo, p.NetworkInfo, p.MountInfo,
		)
		logger.Log.Print(2, "containerid : %s, len(%d)", p.ID, len(p.ID))
	}

	query += strings.Join(placeholders, ", ")

	query += `
		ON DUPLICATE KEY UPDATE
			container_name = VALUES(container_name),
			image          = VALUES(image),
			platform       = VALUES(platform),
			restart_count  = VALUES(restart_count),
			state_info     = VALUES(state_info),
			config_info    = VALUES(config_info),
			network_info   = VALUES(network_info),
			mount_info     = VALUES(mount_info),
			changed_at     = VALUES(changed_at)`

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		logger.Log.Error("failed to upsert container inspect: %v", err)
		return fmt.Errorf("failed to upsert container inspect: %w", err)
	}

	if err = tx.Commit(); err != nil {
		logger.Log.Error("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
