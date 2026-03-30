package service

import (
	"context"

	"agent-service/internal/db"
	"agent-service/internal/dto"
)

type ServiceInterface interface {
	Test()
	CreateContainerInfo(ctx context.Context, req dto.ContainerListData, agentid, hostid int) error
	CreateContainerInspect(ctx context.Context, req dto.ContainerInspectData, agentid, hostid int) error
	CreateContainerStats(ctx context.Context, req dto.ContainerStatsData, agentid, hostid int) error
	CreateContainerEvent(ctx context.Context, req dto.ContainerEvent, agentid, hostid int) error

	ReadHost(ctx context.Context, agentid int) ([]db.Host, error)
	ReadContainerInfo(ctx context.Context, agentid, hostid int) ([]db.ContainerInfo, error)
	ReadContainerInspect(ctx context.Context, agentid, hostid int, containerID string) (*db.ContainerInspect, error)
	ReadContainerStats(ctx context.Context, agentid, hostid int) ([]db.ContainerStats, error)
}
