package api

import (
	"net/http"
	"time"

	"agent-service/internal/logger"

	"github.com/gin-gonic/gin"
)

func (server *Server) testapi(ctx *gin.Context) {
	time.Sleep(time.Microsecond * 3000)

	strdt, err := server.dbHnd.ReadSysdate(ctx)
	if err != nil {
		logger.Log.Error("testapi err..%v", err)
	}
	logger.Log.Print(2, "testapi :%v", strdt)

	ctx.JSON(http.StatusOK, "hello")
}

func (server *Server) heartbeat(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

func (server *Server) terminate(ctx *gin.Context) {
	server.ch_terminate <- true
	logger.Log.Print(2, "Accept terminate command..")
	ctx.JSON(http.StatusOK, nil)
}

func (server *Server) dockerPs2(ctx *gin.Context) {
	var req requestAgentHost
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	containers, err := server.service.ReadContainerInfo(ctx, req.AgentId, req.HostId)
	if err != nil {
		logger.Log.Error("Service Container list error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	response := ToContainerListResponse(containers)
	ctx.JSON(http.StatusOK, SuccessResponse(response))
}

func (server *Server) containerInspect2(ctx *gin.Context) {
	var req requestAgentHostContainer
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	inspect, err := server.service.ReadContainerInspect(ctx, req.AgentId, req.HostId, req.ContainerID)
	if err != nil {
		logger.Log.Error("Service ContainerInspect error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(ToContainerInspectResponse(*inspect)))
}

func (server *Server) statContainer3(ctx *gin.Context) {
	var req requestAgentHostOnly
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	stats, err := server.service.ReadContainerStats(ctx, req.AgentId, req.HostId)
	if err != nil {
		logger.Log.Error("Service ContainerStats error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(ToContainerStatsListResponse(stats)))
}
