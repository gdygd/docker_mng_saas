package db

import (
	"context"
	"database/sql"
)

type DbHandler interface {
	Init() error
	Close(*sql.DB)
	ReadSysdate(ctx context.Context) (string, error)

	CreateContainerInfo(ctx context.Context, params []ContainerInfoParams) error
	UpsertContainerInspect(ctx context.Context, params []ContainerInspectParams) error
	InsertContainerStats(ctx context.Context, params []ContainerStatsParams) error
	InsertContainerEvent(ctx context.Context, param ContainerEventParams) error

	ReadHost(ctx context.Context, agentid int) ([]Host, error)
	ReadContainerInfo(ctx context.Context, agentid, hostid int) ([]ContainerInfo, error)
	ReadContainerInspect(ctx context.Context, agentid, hostid int, containerID string) (*ContainerInspect, error)
	ReadContainerStats(ctx context.Context, agentid, hostid int) ([]ContainerStats, error)
}
