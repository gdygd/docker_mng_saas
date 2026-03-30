package service

import (
	"context"
	"encoding/json"
	"fmt"

	"agent-service/internal/db"
	"agent-service/internal/dto"
	"agent-service/internal/logger"
	"agent-service/internal/memory"
	"agent-service/internal/service"
)

type RpcService struct {
	dbHnd db.DbHandler
	objdb *memory.RedisDb
}

func NewRpcService(dbHnd db.DbHandler, objdb *memory.RedisDb) service.ServiceInterface {
	return &RpcService{
		dbHnd: dbHnd,
		objdb: objdb,
	}
}

func (s *RpcService) Test() {
	fmt.Printf("test service")
}

func (s *RpcService) CreateContainerInfo(ctx context.Context, req dto.ContainerListData, agentid, hostid int) error {
	logger.Log.Print(2, "CreateContainerInfo service...")
	params := toContainerInfoParams(req)
	return s.dbHnd.CreateContainerInfo(ctx, agentid, hostid, params)
}

func (s *RpcService) CreateContainerInspect(ctx context.Context, req dto.ContainerInspectData, agentid, hostid int) error {
	params, err := toContainerInspectParams(req)
	if err != nil {
		return fmt.Errorf("toContainerInspectParams: %w", err)
	}
	return s.dbHnd.UpsertContainerInspect(ctx, agentid, hostid, params)
}

func (s *RpcService) CreateContainerStats(ctx context.Context, req dto.ContainerStatsData, agentid, hostid int) error {
	params := toContainerStatsParams(req)
	return s.dbHnd.InsertContainerStats(ctx, agentid, hostid, params)
}

func (s *RpcService) CreateContainerEvent(ctx context.Context, req dto.ContainerEvent, agentid, hostid int) error {
	param, err := toContainerEventParam(req)
	if err != nil {
		return fmt.Errorf("toContainerEventParam: %w", err)
	}
	return s.dbHnd.InsertContainerEvent(ctx, agentid, hostid, param)
}

func (s *RpcService) ReadHost(ctx context.Context, agentid int) ([]db.Host, error) {
	return nil, nil
}

func (s *RpcService) ReadContainerInfo(ctx context.Context, agentid, hostid int) ([]db.ContainerInfo, error) {
	return nil, nil
}

func (s *RpcService) ReadContainerInspect(ctx context.Context, agentid, hostid int, containerID string) (*db.ContainerInspect, error) {
	return nil, nil
}

func (s *RpcService) ReadContainerStats(ctx context.Context, agentid, hostid int) ([]db.ContainerStats, error) {
	return nil, nil
}

// --- converter ---

func toContainerInfoParams(req dto.ContainerListData) []db.ContainerInfoParams {
	params := make([]db.ContainerInfoParams, 0, len(req.Containers))
	for _, r := range req.Containers {
		params = append(params, db.ContainerInfoParams{
			ID:     r.ID,
			Name:   r.Name,
			Image:  r.Image,
			State:  r.State,
			Status: r.Status,
		})
	}
	return params
}

func toContainerInspectParams(req dto.ContainerInspectData) ([]db.ContainerInspectParams, error) {
	params := make([]db.ContainerInspectParams, 0, len(req.Inspects))
	for _, r := range req.Inspects {
		stateJSON, err := marshalJSON(r.State)
		if err != nil {
			return nil, fmt.Errorf("state marshal: %w", err)
		}
		configJSON, err := marshalJSON(r.Config)
		if err != nil {
			return nil, fmt.Errorf("config marshal: %w", err)
		}
		networkJSON, err := marshalJSON(r.Network)
		if err != nil {
			return nil, fmt.Errorf("network marshal: %w", err)
		}
		mountJSON, err := marshalJSON(r.Mounts)
		if err != nil {
			return nil, fmt.Errorf("mount marshal: %w", err)
		}
		params = append(params, db.ContainerInspectParams{
			ID:           r.ID,
			Name:         r.Name,
			Image:        r.Image,
			Platform:     r.Platform,
			RestartCount: r.RestartCount,
			StateInfo:    stateJSON,
			ConfigInfo:   configJSON,
			NetworkInfo:  networkJSON,
			MountInfo:    mountJSON,
		})
	}
	return params, nil
}

func toContainerStatsParams(req dto.ContainerStatsData) []db.ContainerStatsParams {
	params := make([]db.ContainerStatsParams, 0, len(req.Stats))
	for _, r := range req.Stats {
		params = append(params, db.ContainerStatsParams{
			ID:            r.ID,
			Name:          r.Name,
			CPUPercent:    r.CPUPercent,
			MemoryUsage:   r.MemoryUsage,
			MemoryLimit:   r.MemoryLimit,
			MemoryPercent: r.MemoryPercent,
			NetworkRx:     r.NetworkRx,
			NetworkTx:     r.NetworkTx,
		})
	}
	return params
}

func toContainerEventParam(req dto.ContainerEvent) (db.ContainerEventParams, error) {
	attrsJSON, err := marshalJSON(req.Attrs)
	if err != nil {
		return db.ContainerEventParams{}, fmt.Errorf("attrs marshal: %w", err)
	}
	return db.ContainerEventParams{
		ContainerID: req.ActorID,
		Hostname:    req.Host,
		Type:        req.Type,
		Action:      req.Action,
		ActorID:     req.ActorID,
		ActorName:   req.ActorName,
		EventTime:   req.Timestamp,
		Attrs:       attrsJSON,
	}, nil
}

func marshalJSON(v any) (string, error) {
	if v == nil {
		return "null", nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
