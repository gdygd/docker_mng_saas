package db

import "database/sql"

type ContainerInfoParams struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	State  string `json:"state"`
	Status string `json:"status"`
}

// ContainerInspectParams — JSON 컬럼은 service 레이어에서 직렬화된 문자열로 전달
type ContainerInspectParams struct {
	ID           string
	Name         string
	Image        string
	Platform     string
	RestartCount int
	StateInfo    string // JSON
	ConfigInfo   string // JSON
	NetworkInfo  string // JSON
	MountInfo    string // JSON
}

type ContainerStatsParams struct {
	ID            string
	Name          string
	CPUPercent    float64
	MemoryUsage   uint64
	MemoryLimit   uint64
	MemoryPercent float64
	NetworkRx     uint64
	NetworkTx     uint64
}

type ContainerEventParams struct {
	ContainerID string
	Hostname    string
	Type        string
	Action      string
	ActorID     string
	ActorName   string
	EventTime   int64
	Attrs       string // JSON
}

type ContainerInfo struct {
	AgentId   int
	HostId    int
	HostName  string
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Image     string       `json:"image"`
	State     string       `json:"state"`
	Status    string       `json:"status"`
	ChangedAt sql.NullTime `json:"changed_at"`
}

type ContainerInspect struct {
	// 기본 정보
	AgentId  int
	HostId   int
	HostName string

	ID           string `json:"id"`
	Name         string `json:"name"`
	Image        string `json:"image"`
	Created      string `json:"created"`
	Platform     string `json:"platform"`
	RestartCount int    `json:"restart_count"`

	// 상태 정보
	State []byte `json:"state,omitempty"`

	// 설정 정보
	Config []byte `json:"config,omitempty"`

	// 네트워크 정보
	Network []byte `json:"network,omitempty"`

	// 마운트 정보
	Mounts []byte `json:"mounts,omitempty"`
}

type ContainerStats struct {
	AgentId  int
	HostId   int
	HostName string

	ID            string  `json:"id"`
	Name          string  `json:"name"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryUsage   string  `json:"memory_usage"` // "1.2 GiB"
	MemoryLimit   string  `json:"memory_limit"` // "4.0 GiB"
	MemoryPercent float64 `json:"memory_percent"`
	NetworkRx     string  `json:"network_rx"` // "1.5 MiB"
	NetworkTx     string  `json:"network_tx"` // "2.3 MiB"

	CollectedAt sql.NullTime `json:"collected_at"`
}
