package gapi

import (
	"context"
	"encoding/json"

	"agent-service/internal/logger"

	pb "agent-service/pb"

	"github.com/gdygd/goglib"
)

const (
	AGENT_ID = 1
	HOST_ID  = 1
)

func (server *Server) ContainerState(ctx context.Context, req *pb.AgentMessage) (*pb.ServerMessage, error) {
	logger.Log.Print(2, "rpc ContainerState")
	logger.Log.Print(2, "type : %v, host : %v", req.GetType(), req.GetHost())
	logger.Log.Print(2, "data : %v", req.GetData())

	agentMsg := parseAgentMessage(req)
	_ = agentMsg

	rsp := &pb.ServerMessage{}

	return rsp, nil
}

func (server *Server) ContainerInfo(ctx context.Context, req *pb.AgentMessage) (*pb.ServerMessage, error) {
	logger.Log.Print(2, "rpc ContainerInfo...")
	logger.Log.Print(2, "type : %v, host : %v", req.GetType(), req.GetHost())
	logger.Log.Print(1, "data : %v", req.GetData())

	agentMsg := parseAgentMessage(req)
	_ = agentMsg

	for i, c := range agentMsg.ListData.Containers {
		logger.Log.Print(2, "(%d) ID: %s, Name : %s, Img : %s, State :%s, Stt: %s",
			i, c.ID, c.Name, c.Image, c.State, c.Status)
	}

	logger.Log.Print(2, "rpc ContainerInfo...2")

	err := server.service.CreateContainerInfo(ctx, agentMsg.ListData, AGENT_ID, HOST_ID)
	if err != nil {
		logger.Log.Error("CreateContainerInfo error.. :%v", err)
	}

	logger.Log.Print(2, "rpc ContainerInfo...3")

	rsp := &pb.ServerMessage{}

	return rsp, nil
}

func (server *Server) ContainerInspect(ctx context.Context, req *pb.AgentMessage) (*pb.ServerMessage, error) {
	logger.Log.Print(1, "rpc ContainerInspect")
	logger.Log.Print(1, "type : %v, host : %v", req.GetType(), req.GetHost())

	agentMsg := parseAgentMessage(req)

	for i, c := range agentMsg.InspectData.Inspects {
		logger.Log.Print(1, "(%d) ID: %s, Name : %s, Img : %s, Created :%s, Platform: %s, restart : %d, Status:%s, host : %s, ip:%s",
			i, c.ID, c.Name, c.Image, c.Created, c.Platform, c.RestartCount, c.State.Status, c.Config.Hostname, c.Network.IPAddress)
	}

	if err := server.service.CreateContainerInspect(ctx, agentMsg.InspectData, AGENT_ID, HOST_ID); err != nil {
		logger.Log.Error("CreateContainerInspect error: %v", err)
	}

	return &pb.ServerMessage{}, nil
}

func (server *Server) ContainerStats(ctx context.Context, req *pb.AgentMessage) (*pb.ServerMessage, error) {
	logger.Log.Print(2, "rpc ContainerStats")
	logger.Log.Print(2, "type : %v, host : %v", req.GetType(), req.GetHost())

	agentMsg := parseAgentMessage(req)

	for i, c := range agentMsg.StatsData.Stats {
		logger.Log.Print(2, "(%d) ID: %s, Name : %s, cpu:%.2f, memU:%d memL : %d, memP:%.2f, rx:%d, tx:%d",
			i, c.ID, c.Name, c.CPUPercent, c.MemoryUsage, c.MemoryLimit, c.MemoryPercent, c.NetworkRx, c.NetworkTx)
	}

	if err := server.service.CreateContainerStats(ctx, agentMsg.StatsData, AGENT_ID, HOST_ID); err != nil {
		logger.Log.Error("CreateContainerStats error: %v", err)
	}

	return &pb.ServerMessage{}, nil
}

func (server *Server) ContainerEvent(ctx context.Context, req *pb.AgentMessage) (*pb.ServerMessage, error) {
	agentMsg := parseAgentMessage(req)

	logger.Log.Print(2, "rpc ContainerEvent type:%s action:%s actor:%s",
		agentMsg.EventData.Type, agentMsg.EventData.Action, agentMsg.EventData.ActorID)

	logger.Log.Print(1, "rcp event host : %s", agentMsg.EventData.Host)
	logger.Log.Print(1, "rcp event type : %s", agentMsg.EventData.Type)
	logger.Log.Print(1, "rcp event action : %s", agentMsg.EventData.Action)
	logger.Log.Print(1, "rcp event actorid : %s", agentMsg.EventData.ActorID)
	logger.Log.Print(1, "rcp event actorname : %s", agentMsg.EventData.ActorName)
	logger.Log.Print(1, "rcp event timestamp : %s", agentMsg.EventData.Timestamp)
	logger.Log.Print(1, "rcp event attrs : %v", agentMsg.EventData.Attrs)

	if err := server.service.CreateContainerEvent(ctx, agentMsg.EventData, AGENT_ID, HOST_ID); err != nil {
		logger.Log.Error("CreateContainerEvent error: %v", err)
		return &pb.ServerMessage{}, nil
	}

	// sse

	data, _ := json.Marshal(agentMsg.EventData)
	logger.Log.Print(2, "sse container event..%s", string(data))
	goglib.SendSSE(goglib.EventData{
		Msgtype: "container-event",
		Data:    string(data),
	})

	return &pb.ServerMessage{}, nil
}
