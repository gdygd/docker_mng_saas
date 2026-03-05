package gapi

import (
	"context"

	"agent-service/internal/logger"
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

func (server *Server) ContainerInfo(ctx context.Context, req *pb.ContainerListData) (*pb.ServerMessage, error) {
	logger.Log.Print(2, "rpc ContainerInfo")

	rsp := &pb.ServerMessage{}

	return rsp, nil
}

func (server *Server) ContainerInspect(ctx context.Context, req *pb.ContainerInspectData) (*pb.ServerMessage, error) {
	logger.Log.Print(2, "rpc ContainerInfo")

	rsp := &pb.ServerMessage{}

	return rsp, nil
}

func (server *Server) ContainerStats(ctx context.Context, req *pb.ContainerStatsData) (*pb.ServerMessage, error) {
	logger.Log.Print(2, "rpc ContainerInfo")

	rsp := &pb.ServerMessage{}

	return rsp, nil
}

func (server *Server) ContainerEvent(ctx context.Context, req *pb.ContainerEventData) (*pb.ServerMessage, error) {
	logger.Log.Print(2, "rpc ContainerInfo")

	rsp := &pb.ServerMessage{}

	return rsp, nil
}
