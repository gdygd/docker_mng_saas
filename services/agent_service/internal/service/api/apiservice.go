package service

import (
	"agent-service/internal/dto"
	"context"
	"fmt"

	"agent-service/internal/db"
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

func (s *ApiService) CreateContainerInfo(ctx context.Context, req dto.ContainerListData, hostid int) error {
	return nil
}

func (s *ApiService) CreateContainerInspect(ctx context.Context, req dto.ContainerInspectData, hostid int) error {
	return nil
}

func (s *ApiService) CreateContainerStats(ctx context.Context, req dto.ContainerStatsData, hostid int) error {
	return nil
}

func (s *ApiService) CreateContainerEvent(ctx context.Context, req dto.ContainerEvent, hostid int) error {
	return nil
}
