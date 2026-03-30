package db

import (
	"context"
	"database/sql"
)

type DbHandler interface {
	Init() error
	Close(*sql.DB)
	ReadSysdate(ctx context.Context) (string, error)

	CreateContainerInfo(ctx context.Context, agentid, hostId int, params []ContainerInfoParams) error
	UpsertContainerInspect(ctx context.Context, agentid, hostId int, params []ContainerInspectParams) error
	InsertContainerStats(ctx context.Context, agentid, hostId int, params []ContainerStatsParams) error
	InsertContainerEvent(ctx context.Context, agentid, hostId int, param ContainerEventParams) error

	ReadHost(ctx context.Context, agentid int) ([]Host, error)
	ReadContainerInfo(ctx context.Context, agentid, hostid int) ([]ContainerInfo, error)
	ReadContainerInspect(ctx context.Context, agentid, hostid int, containerID string) (*ContainerInspect, error)
	ReadContainerStats(ctx context.Context, agentid, hostid int) ([]ContainerStats, error)
}
