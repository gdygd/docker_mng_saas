package gapi

import (
	"context"
	"encoding/json"

	"agent-service/internal/logger"

	pb "agent-service/pb"

	"github.com/gdygd/goglib"
)

// HOST_ID는 현재 테스트용 상수. 실제 운영 시 req에서 파싱 필요.

const (
	AGENT_ID = 1
	HOST_ID  = 1
)

func (server *Server) ContainerState(ctx context.Context, req *pb.AgentMessage) (*pb.ServerMessage, error) {
	logger.Log.Print(1, "rpc ContainerState")
	logger.Log.Print(1, "agentid : %v type : %v, host : %v", req.GetAgentid(), req.GetType(), req.GetHost())
	logger.Log.Print(1, "data : %v", req.GetData())

	agentMsg := parseAgentMessage(req)
	_ = agentMsg

	rsp := &pb.ServerMessage{}

	return rsp, nil
}

func (server *Server) ContainerInfo(ctx context.Context, req *pb.AgentMessage) (*pb.ServerMessage, error) {
	logger.Log.Print(1, "rpc ContainerInfo agent[%d] host:%v", req.GetAgentid(), req.GetHost())

	agentMsg := parseAgentMessage(req)
	server.batch.PushContainerInfo(int(req.GetAgentid()), HOST_ID, agentMsg.ListData)

	return &pb.ServerMessage{}, nil
}

func (server *Server) ContainerInspect(ctx context.Context, req *pb.AgentMessage) (*pb.ServerMessage, error) {
	logger.Log.Print(1, "rpc ContainerInspect agent[%d] host:%v", req.GetAgentid(), req.GetHost())

	agentMsg := parseAgentMessage(req)
	server.batch.PushContainerInspect(int(req.GetAgentid()), HOST_ID, agentMsg.InspectData)

	return &pb.ServerMessage{}, nil
}

func (server *Server) ContainerStats(ctx context.Context, req *pb.AgentMessage) (*pb.ServerMessage, error) {
	// server.statsCounter.inc()

	agentMsg := parseAgentMessage(req)

	server.batch.PushContainerStats(int(req.GetAgentid()), HOST_ID, agentMsg.StatsData)

	return &pb.ServerMessage{}, nil
}

func (server *Server) ContainerEvent(ctx context.Context, req *pb.AgentMessage) (*pb.ServerMessage, error) {
	agentMsg := parseAgentMessage(req)

	logger.Log.Print(1, "rpc ContainerEvent agent[%d] type:%s action:%s actor:%s",
		req.GetAgentid(), agentMsg.EventData.Type, agentMsg.EventData.Action, agentMsg.EventData.ActorID)

	// SSE는 즉시 전송 (DB 쓰기와 무관)
	data, _ := json.Marshal(agentMsg.EventData)
	goglib.SendSSE(goglib.EventData{
		Msgtype: "container-event",
		Data:    string(data),
	})

	// DB 쓰기는 배치 큐에 위임
	server.batch.PushContainerEvent(int(req.GetAgentid()), HOST_ID, agentMsg.EventData)

	return &pb.ServerMessage{}, nil
}
