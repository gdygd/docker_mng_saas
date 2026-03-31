package dbbatch

import "agent-service/internal/dto"

// ContainerInfoItem — gRPC ContainerInfo 수신 데이터
type ContainerInfoItem struct {
	AgentId int
	HostId  int
	Data    dto.ContainerListData
}

// ContainerInspectItem — gRPC ContainerInspect 수신 데이터
type ContainerInspectItem struct {
	AgentId int
	HostId  int
	Data    dto.ContainerInspectData
}

// ContainerStatsItem — gRPC ContainerStats 수신 데이터
type ContainerStatsItem struct {
	AgentId int
	HostId  int
	Data    dto.ContainerStatsData
}

// ContainerEventItem — gRPC ContainerEvent 수신 데이터
type ContainerEventItem struct {
	AgentId int
	HostId  int
	Data    dto.ContainerEvent
}
