package db

import (
	"context"
	"database/sql"
)

type DbHandler interface {
	Init() error
	Close(*sql.DB)
	ReadSysdate(ctx context.Context) (string, error)

	CreateContainerInfo(ctx context.Context, hostId int, params []ContainerInfoParams) error
	UpsertContainerInspect(ctx context.Context, hostId int, params []ContainerInspectParams) error
	InsertContainerStats(ctx context.Context, hostId int, params []ContainerStatsParams) error
	InsertContainerEvent(ctx context.Context, hostId int, param ContainerEventParams) error
}
