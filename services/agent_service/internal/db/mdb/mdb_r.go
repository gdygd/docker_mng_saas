package mdb

import (
	"context"
	"fmt"

	"agent-service/internal/db"
	"agent-service/internal/logger"
)

func (q *MariaDbHandler) ReadSysdate(ctx context.Context) (string, error) {
	ado := q.GetDB()

	query := `
	select now() as dt from dual
	`

	rows, err := ado.QueryContext(ctx, query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	strDateTime := ""
	if rows.Next() {
		if err := rows.Scan(
			&strDateTime,
		); err != nil {
			return "", err
		}
	}
	if err := rows.Close(); err != nil {
		return "", err
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return strDateTime, nil
}

func (q *MariaDbHandler) ReadContainerInfo(ctx context.Context, agentid, hostid int) ([]db.ContainerInfo, error) {
	ado := q.GetDB()

	query := `
	select ah.id
		 , ah.host_id 
		 , ifnull(ah.hostname, '') hostname
		 , ci.container_id  
		 , ifnull(ci.container_name, '') container_name
		 , ifnull(ci.image, '') image
		 , ifnull(ci.state, '') state
		 , ifnull(ci.status, '') status
		 , ci.changed_at
	from agent_host ah
	inner join container_info ci
	on ah.id = ? 
	and ah.host_id  = ?
	`

	rows, err := ado.QueryContext(ctx, query, agentid, hostid)
	if err != nil {
		logger.Log.Error("ReadContainerInfo#1 error %v", err)
		return nil, err
	}
	defer rows.Close()

	var rst []db.ContainerInfo = []db.ContainerInfo{}

	for rows.Next() {
		row := db.ContainerInfo{}
		if err := rows.Scan(
			&row.AgentId,
			&row.HostId,
			&row.HostName,
			&row.ID,
			&row.Name,
			&row.Image,
			&row.State,
			&row.Status,
			&row.ChangedAt,
		); err != nil {
			logger.Log.Error("ReadContainerInfo#2 error %v", err)
			return nil, err
		}
		rst = append(rst, row)
	}
	if err := rows.Close(); err != nil {
		logger.Log.Error("ReadContainerInfo#3 error %v", err)
		return nil, err
	}
	if err := rows.Err(); err != nil {
		logger.Log.Error("ReadContainerInfo#4 error %v", err)
		return nil, err
	}
	return rst, nil
}

func (q *MariaDbHandler) ReadContainerInspect(ctx context.Context, agentid, hostid int, containerID string) (*db.ContainerInspect, error) {
	ado := q.GetDB()

	query := `
	SELECT ci.container_id
		 , IFNULL(ci.container_name, '') container_name
		 , IFNULL(ci.image, '') image
		 , IFNULL(ci.platform, '') platform
		 , IFNULL(ci.restart_count, 0) restart_count
		 , ci.state_info
		 , ci.config_info
		 , ci.network_info
		 , ci.mount_info
		 , IFNULL(DATE_FORMAT(ci.changed_at, '%Y-%m-%dT%H:%i:%sZ'), '') changed_at
	FROM agent_host ah
	INNER JOIN container_inspect ci ON ci.host_id = ah.host_id
	WHERE ah.id = ?
	AND ah.host_id = ?
	AND ci.container_id LIKE CONCAT(?, '%')
	`

	row := ado.QueryRowContext(ctx, query, agentid, hostid, containerID)

	r := &db.ContainerInspect{AgentId: agentid, HostId: hostid}
	if err := row.Scan(
		&r.ID,
		&r.Name,
		&r.Image,
		&r.Platform,
		&r.RestartCount,
		&r.State,
		&r.Config,
		&r.Network,
		&r.Mounts,
		&r.Created,
	); err != nil {
		logger.Log.Error("ReadContainerInspect error %v", err)
		return nil, err
	}
	return r, nil
}

func (q *MariaDbHandler) ReadContainerStats(ctx context.Context, agentid, hostid int) ([]db.ContainerStats, error) {
	ado := q.GetDB()

	query := `
	SELECT cs.container_id
		, IFNULL(cs.container_name, '') container_name
		, IFNULL(cs.cpu_percent, 0) cpu_percent
		, IFNULL(cs.memory_usage, 0) memory_usage
		, IFNULL(cs.memory_limit, 0) memory_limit
		, IFNULL(cs.memory_percent, 0) memory_percent
		, IFNULL(cs.network_rx, 0) network_rx
		, IFNULL(cs.network_tx, 0) network_tx
		, cs.collected_at
	FROM agent_host ah
	INNER JOIN container_stats cs ON  cs.id = ah.id AND cs.host_id = ah.host_id
	WHERE ah.id = ?
	and ah.HOST_ID = ?
	`

	rows, err := ado.QueryContext(ctx, query, agentid, hostid)
	if err != nil {
		logger.Log.Error("ReadContainerStats#1 error %v", err)
		return nil, err
	}
	defer rows.Close()

	var rst []db.ContainerStats = []db.ContainerStats{}

	for rows.Next() {
		r := db.ContainerStats{AgentId: agentid, HostId: hostid}
		var memUsage, memLimit, netRx, netTx uint64
		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.CPUPercent,
			&memUsage,
			&memLimit,
			&r.MemoryPercent,
			&netRx,
			&netTx,
			&r.CollectedAt,
		); err != nil {
			logger.Log.Error("ReadContainerStats#2 error %v", err)
			return nil, err
		}
		r.MemoryUsage = formatBytes(memUsage)
		r.MemoryLimit = formatBytes(memLimit)
		r.NetworkRx = formatBytes(netRx)
		r.NetworkTx = formatBytes(netTx)
		rst = append(rst, r)
	}
	if err := rows.Close(); err != nil {
		logger.Log.Error("ReadContainerStats#3 error %v", err)
		return nil, err
	}
	if err := rows.Err(); err != nil {
		logger.Log.Error("ReadContainerStats#4 error %v", err)
		return nil, err
	}
	return rst, nil
}

func formatBytes(n uint64) string {
	const (
		GiB = 1024 * 1024 * 1024
		MiB = 1024 * 1024
		KiB = 1024
	)
	switch {
	case n >= GiB:
		return fmt.Sprintf("%.1f GiB", float64(n)/GiB)
	case n >= MiB:
		return fmt.Sprintf("%.1f MiB", float64(n)/MiB)
	case n >= KiB:
		return fmt.Sprintf("%.1f KiB", float64(n)/KiB)
	default:
		return fmt.Sprintf("%d B", n)
	}
}
