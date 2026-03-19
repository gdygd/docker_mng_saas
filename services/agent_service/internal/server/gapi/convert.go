package gapi

import (
	"agent-service/internal/dto"
	pb "agent-service/pb"
	"time"
)

func parseAgentMessage(req *pb.AgentMessage) dto.AgentMessage {
	msg := dto.AgentMessage{
		AgentKey:  req.GetAgentKey(),
		DataType:  int(req.GetType()),
		Host:      req.GetHost(),
		Timestamp: time.Unix(req.GetTimestamp(), 0),
	}

	switch req.GetType() {
	case pb.DataType_CONTAINER_LIST:
		msg.ListData = parseContainerListData(req.GetListData())
	case pb.DataType_CONTAINER_INSPECT:
		msg.InspectData = parseContainerInspectData(req.GetInspectData())
	case pb.DataType_CONTAINER_STATS:
		msg.StatsData = parseContainerStatsData(req.GetStatsData())
	case pb.DataType_CONTAINER_EVENT:
		msg.EventData = parseContainerEventData(req.GetHost(), req.GetEventData())
	}

	return msg
}

func parseContainerListData(src *pb.ContainerListData) dto.ContainerListData {
	if src == nil {
		return dto.ContainerListData{}
	}
	containers := make([]dto.ContainerInfo, len(src.GetContainers()))
	for i, c := range src.GetContainers() {
		containers[i] = dto.ContainerInfo{
			ID:     c.GetId(),
			Name:   c.GetName(),
			Image:  c.GetImage(),
			State:  c.GetState(),
			Status: c.GetStatus(),
		}
	}
	return dto.ContainerListData{Containers: containers}
}

func parseContainerInspectData(src *pb.ContainerInspectData) dto.ContainerInspectData {
	if src == nil {
		return dto.ContainerInspectData{}
	}
	inspects := make([]dto.ContainerInspectInfo, len(src.GetInspects()))
	for i, ins := range src.GetInspects() {
		info := dto.ContainerInspectInfo{
			ID:       ins.GetId(),
			Name:     ins.GetName(),
			Image:    ins.GetImage(),
			Created:  ins.GetCreated(),
			Platform: ins.GetPlatform(),
		}
		if s := ins.GetState(); s != nil {
			info.State = &dto.ContainerStateInfo{
				Status:     s.GetStatus(),
				Running:    s.GetRunning(),
				Paused:     s.GetPaused(),
				Restarting: s.GetRestarting(),
				ExitCode:   int(s.GetExitCode()),
				StartedAt:  s.GetStartedAt(),
				FinishedAt: s.GetFinishedAt(),
			}
		}
		if c := ins.GetConfig(); c != nil {
			info.Config = &dto.ContainerConfigInfo{
				Hostname:   c.GetHostname(),
				User:       c.GetUser(),
				Env:        c.GetEnv(),
				Cmd:        c.GetCmd(),
				Entrypoint: c.GetEntrypoint(),
				WorkingDir: c.GetWorkingDir(),
				Labels:     c.GetLabels(),
			}
		}
		if n := ins.GetNetwork(); n != nil {
			ports := make(map[string][]dto.PortBindingInfo)
			for port, bindings := range n.GetPorts() {
				var bindList []dto.PortBindingInfo
				for _, b := range bindings.GetBindings() {
					bindList = append(bindList, dto.PortBindingInfo{
						HostIP:   b.GetHostIp(),
						HostPort: b.GetHostPort(),
					})
				}
				ports[port] = bindList
			}
			networks := make(map[string]dto.NetworkEndpoint)
			for name, ep := range n.GetNetworks() {
				networks[name] = dto.NetworkEndpoint{
					NetworkID:  ep.GetNetworkId(),
					IPAddress:  ep.GetIpAddress(),
					Gateway:    ep.GetGateway(),
					MacAddress: ep.GetMacAddress(),
				}
			}
			info.Network = &dto.ContainerNetworkInfo{
				IPAddress:  n.GetIpAddress(),
				Gateway:    n.GetGateway(),
				MacAddress: n.GetMacAddress(),
				Ports:      ports,
				Networks:   networks,
			}
		}
		mounts := make([]dto.MountPointInfo, len(ins.GetMounts()))
		for j, m := range ins.GetMounts() {
			mounts[j] = dto.MountPointInfo{
				Type:        m.GetType(),
				Name:        m.GetName(),
				Source:      m.GetSource(),
				Destination: m.GetDestination(),
				Mode:        m.GetMode(),
				RW:          m.GetRw(),
			}
		}
		info.Mounts = mounts
		inspects[i] = info
	}
	return dto.ContainerInspectData{Inspects: inspects}
}

func parseContainerStatsData(src *pb.ContainerStatsData) dto.ContainerStatsData {
	if src == nil {
		return dto.ContainerStatsData{}
	}
	stats := make([]dto.ContainerStatsInfo, len(src.GetStats()))
	for i, s := range src.GetStats() {
		stats[i] = dto.ContainerStatsInfo{
			ID:            s.GetId(),
			Name:          s.GetName(),
			CPUPercent:    s.GetCpuPercent(),
			MemoryUsage:   s.GetMemoryUsage(),
			MemoryLimit:   s.GetMemoryLimit(),
			MemoryPercent: s.GetMemoryPercent(),
			NetworkRx:     s.GetNetworkRx(),
			NetworkTx:     s.GetNetworkTx(),
		}
	}
	return dto.ContainerStatsData{Stats: stats}
}

func parseContainerEventData(host string, src *pb.ContainerEventData) dto.ContainerEvent {
	if src == nil {
		return dto.ContainerEvent{}
	}
	return dto.ContainerEvent{
		Host:      host,
		Type:      src.GetType(),
		Action:    src.GetAction(),
		ActorID:   src.GetActorId(),
		ActorName: src.GetActorName(),
		Timestamp: src.GetTimestamp(),
		Attrs:     src.GetAttrs(),
	}
}
