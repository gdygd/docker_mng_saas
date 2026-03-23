package api

type requestAgentHost struct {
	AgentId int `uri:"agentid" binding:"required"`
	HostId  int `uri:"hostid" binding:"required"`
}

type requestAgentHostOnly struct {
	AgentId int `uri:"agentid" binding:"required"`
	HostId  int `uri:"host" binding:"required"`
}

type requestAgentHostContainer struct {
	AgentId     int    `uri:"agentid" binding:"required"`
	HostId      int    `uri:"host" binding:"required"`
	ContainerID string `uri:"id" binding:"required"`
}
