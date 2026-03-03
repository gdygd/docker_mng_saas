package gapi

import (
	"context"

	"agent-service/internal/logger"
	pb "agent-service/pb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	logger.Log.Print(2, "rpc loginuser")
	logger.Log.Print(2, "user : %s, pw : %s", req.GetUsername(), req.GetPassword())

	rsp := &pb.LoginUserResponse{}

	return rsp, nil
}
