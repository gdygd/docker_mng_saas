package service

import (
	"context"
	"fmt"

	"agent-service/internal/db"
	"agent-service/internal/dto"
	"agent-service/internal/memory"
	"agent-service/internal/service"
)

type ApiService struct {
	dbHnd db.DbHandler
	objdb *memory.RedisDb
}

func NewApiService(dbHnd db.DbHandler, objdb *memory.RedisDb) service.ServiceInterface {
	return &ApiService{
		dbHnd: dbHnd,
		objdb: objdb,
	}
}

func (s *ApiService) Test() {
	fmt.Printf("test service")
}

func (s *ApiService) CreateContainerInfo(ctx context.Context, req dto.ContainerListData, agentid, hostid int) error {
	return nil
}

func (s *ApiService) CreateContainerInspect(ctx context.Context, req dto.ContainerInspectData, agentid, hostid int) error {
	return nil
}

func (s *ApiService) CreateContainerStats(ctx context.Context, req dto.ContainerStatsData, agentid, hostid int) error {
	return nil
}

func (s *ApiService) CreateContainerEvent(ctx context.Context, req dto.ContainerEvent, agentid, hostid int) error {
	return nil
}

func (s *ApiService) ReadHost(ctx context.Context, agentid int) ([]db.Host, error) {
	return s.dbHnd.ReadHost(ctx, agentid)
}

func (s *ApiService) ReadContainerInfo(ctx context.Context, agentid, hostid int) ([]db.ContainerInfo, error) {
	return s.dbHnd.ReadContainerInfo(ctx, agentid, hostid)
}

func (s *ApiService) ReadContainerInspect(ctx context.Context, agentid, hostid int, containerID string) (*db.ContainerInspect, error) {
	return s.dbHnd.ReadContainerInspect(ctx, agentid, hostid, containerID)
}

func (s *ApiService) ReadContainerStats(ctx context.Context, agentid, hostid int) ([]db.ContainerStats, error) {
	return s.dbHnd.ReadContainerStats(ctx, agentid, hostid)
}
