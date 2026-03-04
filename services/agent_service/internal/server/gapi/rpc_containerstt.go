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

	rsp := &pb.ServerMessage{}

	return rsp, nil
}
