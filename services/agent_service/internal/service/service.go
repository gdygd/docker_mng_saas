package service

import (
	"agent-service/internal/dto"
	"context"
)

type ServiceInterface interface {
	Test()
	CreateContainerInfo(ctx context.Context, req dto.ContainerListData, hostid int) error
	CreateContainerInspect(ctx context.Context, req dto.ContainerInspectData, hostid int) error
	CreateContainerStats(ctx context.Context, req dto.ContainerStatsData, hostid int) error
	CreateContainerEvent(ctx context.Context, req dto.ContainerEvent, hostid int) error
}
