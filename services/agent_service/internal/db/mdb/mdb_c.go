package mdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"agent-service/internal/db"
	"agent-service/internal/logger"
)

func (q *MariaDbHandler) CreateContainerInfo(ctx context.Context, params []db.ContainerInfoParams) error {
	logger.Log.Print(1, "CreateContainerInfo db.. len(%d)", len(params))
	if len(params) == 0 {
		return nil
	}

	ado := q.GetDB()

	// 트랜잭션 시작
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

	// Bulk INSERT 쿼리 빌드 (placeholder 방식)
	query := `
		INSERT INTO container_info (
			id,
			host_id,
			container_id,
			container_name,
			image,
			state,
			status,
			changed_at
		) VALUES `

	placeholders := make([]string, 0, len(params))
	args := make([]interface{}, 0, len(params)*7)

	for _, arg := range params {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, NOW())")
		args = append(args, arg.AgentId, arg.HostId, arg.ID, arg.Name, arg.Image, arg.State, arg.Status)
		logger.Log.Print(1, "agent[%d], host[%d] containerID:%s, len(%d)", arg.AgentId, arg.HostId, arg.ID, len(arg.ID))
	}

	query += strings.Join(placeholders, ", ")

	query += `
		ON DUPLICATE KEY UPDATE
			container_name = VALUES(container_name),
			image          = VALUES(image),
			state          = VALUES(state),
			status         = VALUES(status),
			changed_at     = VALUES(changed_at)`

	logger.Log.Print(1, "CreateContainerInfo qry :%s", query)

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		logger.Log.Error("failed to insert container info: %v", err)
		return fmt.Errorf("failed to insert container info: %w", err)
	}

	if err = tx.Commit(); err != nil {
		logger.Log.Error("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
