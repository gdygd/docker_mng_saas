package gapi

import (
	pb "agent-service/pb"
	"time"
)

func parseAgentMessage(req *pb.AgentMessage) AgentMessage {
	msg := AgentMessage{
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

func parseContainerListData(src *pb.ContainerListData) ContainerListData {
	if src == nil {
		return ContainerListData{}
	}
	containers := make([]ContainerInfo, len(src.GetContainers()))
	for i, c := range src.GetContainers() {
		containers[i] = ContainerInfo{
			ID:     c.GetId(),
			Name:   c.GetName(),
			Image:  c.GetImage(),
			State:  c.GetState(),
			Status: c.GetStatus(),
		}
	}
	return ContainerListData{Containers: containers}
}

func parseContainerInspectData(src *pb.ContainerInspectData) ContainerInspectData {
	if src == nil {
		return ContainerInspectData{}
	}
	inspects := make([]ContainerInspectInfo, len(src.GetInspects()))
	for i, ins := range src.GetInspects() {
		info := ContainerInspectInfo{
			ID:       ins.GetId(),
			Name:     ins.GetName(),
			Image:    ins.GetImage(),
			Created:  ins.GetCreated(),
			Platform: ins.GetPlatform(),
		}
		if s := ins.GetState(); s != nil {
			info.State = &ContainerStateInfo{
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
			info.Config = &ContainerConfigInfo{
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
			ports := make(map[string][]PortBindingInfo)
			for port, bindings := range n.GetPorts() {
				var bindList []PortBindingInfo
				for _, b := range bindings.GetBindings() {
					bindList = append(bindList, PortBindingInfo{
						HostIP:   b.GetHostIp(),
						HostPort: b.GetHostPort(),
					})
				}
				ports[port] = bindList
			}
			networks := make(map[string]NetworkEndpoint)
			for name, ep := range n.GetNetworks() {
				networks[name] = NetworkEndpoint{
					NetworkID:  ep.GetNetworkId(),
					IPAddress:  ep.GetIpAddress(),
					Gateway:    ep.GetGateway(),
					MacAddress: ep.GetMacAddress(),
				}
			}
			info.Network = &ContainerNetworkInfo{
				IPAddress:  n.GetIpAddress(),
				Gateway:    n.GetGateway(),
				MacAddress: n.GetMacAddress(),
				Ports:      ports,
				Networks:   networks,
			}
		}
		mounts := make([]MountPointInfo, len(ins.GetMounts()))
		for j, m := range ins.GetMounts() {
			mounts[j] = MountPointInfo{
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
	return ContainerInspectData{Inspects: inspects}
}

func parseContainerStatsData(src *pb.ContainerStatsData) ContainerStatsData {
	if src == nil {
		return ContainerStatsData{}
	}
	stats := make([]ContainerStatsInfo, len(src.GetStats()))
	for i, s := range src.GetStats() {
		stats[i] = ContainerStatsInfo{
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
	return ContainerStatsData{Stats: stats}
}

func parseContainerEventData(host string, src *pb.ContainerEventData) ContainerEvent {
	if src == nil {
		return ContainerEvent{}
	}
	return ContainerEvent{
		Host:      host,
		Type:      src.GetType(),
		Action:    src.GetAction(),
		ActorID:   src.GetActorId(),
		ActorName: src.GetActorName(),
		Timestamp: src.GetTimestamp(),
		Attrs:     src.GetAttrs(),
	}
}
