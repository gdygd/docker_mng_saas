package gapi

import (
	"agent-service/internal/logger"
	"context"

	pb "agent-service/pb"
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
	logger.Log.Print(1, "rpc ContainerInfo")
	logger.Log.Print(1, "type : %v, host : %v", req.GetType(), req.GetHost())
	logger.Log.Print(1, "data : %v", req.GetData())

	agentMsg := parseAgentMessage(req)
	_ = agentMsg

	for i, c := range agentMsg.ListData.Containers {
		logger.Log.Print(1, "(%d) ID: %s, Name : %s, Img : %s, State :%s, Stt: %s",
			i, c.ID, c.Name, c.Image, c.State, c.Status)
	}

	rsp := &pb.ServerMessage{}

	return rsp, nil
}

func (server *Server) ContainerInspect(ctx context.Context, req *pb.AgentMessage) (*pb.ServerMessage, error) {
	logger.Log.Print(1, "rpc ContainerInspect")
	logger.Log.Print(1, "type : %v, host : %v", req.GetType(), req.GetHost())
	logger.Log.Print(1, "data : %v", req.GetData())

	agentMsg := parseAgentMessage(req)
	_ = agentMsg

	for i, c := range agentMsg.InspectData.Inspects {
		logger.Log.Print(1, "(%d) ID: %s, Name : %s, Img : %s, Created :%s, Platform: %s, restart : %d, Status:%s, host : %s, ip:%s",
			i, c.ID, c.Name, c.Image, c.Created, c.Platform, c.RestartCount, c.State.Status, c.Config.Hostname, c.Network.IPAddress)
	}

	rsp := &pb.ServerMessage{}

	return rsp, nil
}

func (server *Server) ContainerStats(ctx context.Context, req *pb.AgentMessage) (*pb.ServerMessage, error) {
	logger.Log.Print(2, "rpc ContainerStats")
	logger.Log.Print(2, "type : %v, host : %v", req.GetType(), req.GetHost())
	logger.Log.Print(2, "data : %v", req.GetData())

	agentMsg := parseAgentMessage(req)
	_ = agentMsg

	for i, c := range agentMsg.StatsData.Stats {
		logger.Log.Print(2, "(%d) ID: %s, Name : %s, cpu:%.2f, memU:%d memL : %d, memP:%.2f, rx:%d, tx:%d",
			i, c.ID, c.Name, c.CPUPercent, c.MemoryUsage, c.MemoryLimit, c.MemoryPercent, c.NetworkRx, c.NetworkTx)
	}

	rsp := &pb.ServerMessage{}

	return rsp, nil
}

func (server *Server) ContainerEvent(ctx context.Context, req *pb.AgentMessage) (*pb.ServerMessage, error) {
	// logger.Log.Print(2, "rpc ContainerEvent")

	// logger.Log.Print(2, "type : %v, host : %v", req.GetType(), req.GetHost())
	// logger.Log.Print(2, "data : %v", req.GetData())

	agentMsg := parseAgentMessage(req)
	_ = agentMsg

	rsp := &pb.ServerMessage{}

	return rsp, nil
}
