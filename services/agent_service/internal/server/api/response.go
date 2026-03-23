package api

import (
	"encoding/json"
	"time"

	"agent-service/internal/db"
)

// ============================================================================
// Generic API Response
// ============================================================================

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func SuccessResponse[T any](data T) APIResponse[T] {
	return APIResponse[T]{
		Success: true,
		Data:    data,
	}
}

func SuccessMessageResponse[T any](message string, data T) APIResponse[T] {
	return APIResponse[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string) APIResponse[any] {
	return APIResponse[any]{
		Success: false,
		Message: message,
	}
}

type ContainerResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	State     string    `json:"state"`
	Status    string    `json:"status"`
	ChangedAt time.Time `json:"changed_at"`
}

func ToContainerResponse(c db.ContainerInfo) ContainerResponse {
	return ContainerResponse{
		ID:        c.ID,
		Name:      c.Name,
		Image:     c.Image,
		State:     c.State,
		Status:    c.Status,
		ChangedAt: c.ChangedAt.Time,
	}
}

func ToContainerListResponse(containers []db.ContainerInfo) []ContainerResponse {
	result := make([]ContainerResponse, 0, len(containers))
	for _, c := range containers {
		result = append(result, ToContainerResponse(c))
	}
	return result
}

// ============================================================================
// Container Inspect Response
// ============================================================================

type ContainerInspectResponse struct {
	// 기본 정보
	ID           string `json:"id"`
	Name         string `json:"name"`
	Image        string `json:"image"`
	Created      string `json:"created"`
	Platform     string `json:"platform"`
	RestartCount int    `json:"restart_count"`

	// 상태 정보
	State *StateResponse `json:"state,omitempty"`

	// 설정 정보
	Config *ConfigResponse `json:"config,omitempty"`

	// 네트워크 정보
	Network *NetworkResponse `json:"network,omitempty"`

	// 마운트 정보
	Mounts []MountResponse `json:"mounts,omitempty"`
}

type StateResponse struct {
	Status     string `json:"status"`
	Running    bool   `json:"running"`
	Paused     bool   `json:"paused"`
	Restarting bool   `json:"restarting"`
	ExitCode   int    `json:"exit_code"`
	StartedAt  string `json:"started_at,omitempty"`
	FinishedAt string `json:"finished_at,omitempty"`
}

type ConfigResponse struct {
	Hostname   string            `json:"hostname,omitempty"`
	User       string            `json:"user,omitempty"`
	Env        []string          `json:"env,omitempty"`
	Cmd        []string          `json:"cmd,omitempty"`
	Entrypoint []string          `json:"entrypoint,omitempty"`
	WorkingDir string            `json:"working_dir,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
}

type NetworkResponse struct {
	IPAddress  string                             `json:"ip_address"`
	Gateway    string                             `json:"gateway"`
	MacAddress string                             `json:"mac_address"`
	Ports      map[string][]PortResponse          `json:"ports,omitempty"`
	Networks   map[string]NetworkEndpointResponse `json:"networks,omitempty"`
}

type PortResponse struct {
	HostIP   string `json:"host_ip"`
	HostPort string `json:"host_port"`
}

type NetworkEndpointResponse struct {
	NetworkID  string `json:"network_id"`
	IPAddress  string `json:"ip_address"`
	Gateway    string `json:"gateway"`
	MacAddress string `json:"mac_address"`
}

type MountResponse struct {
	Type        string `json:"type"`
	Name        string `json:"name,omitempty"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	ReadWrite   bool   `json:"rw"`
}

func ToContainerInspectResponse(c db.ContainerInspect) ContainerInspectResponse {
	resp := ContainerInspectResponse{
		ID:           c.ID,
		Name:         c.Name,
		Image:        c.Image,
		Created:      c.Created,
		Platform:     c.Platform,
		RestartCount: c.RestartCount,
	}
	if len(c.State) > 0 {
		var s StateResponse
		if err := json.Unmarshal(c.State, &s); err == nil {
			resp.State = &s
		}
	}
	if len(c.Config) > 0 {
		var cfg ConfigResponse
		if err := json.Unmarshal(c.Config, &cfg); err == nil {
			resp.Config = &cfg
		}
	}
	if len(c.Network) > 0 {
		var n NetworkResponse
		if err := json.Unmarshal(c.Network, &n); err == nil {
			resp.Network = &n
		}
	}
	if len(c.Mounts) > 0 {
		var m []MountResponse
		if err := json.Unmarshal(c.Mounts, &m); err == nil {
			resp.Mounts = m
		}
	}
	return resp
}

// ============================================================================
// Container Stats Response
// ============================================================================

type ContainerStatsResponse struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryUsage   string  `json:"memory_usage"` // "1.2 GiB"
	MemoryLimit   string  `json:"memory_limit"` // "4.0 GiB"
	MemoryPercent float64 `json:"memory_percent"`
	NetworkRx     string  `json:"network_rx"` // "1.5 MiB"
	NetworkTx     string  `json:"network_tx"` // "2.3 MiB"
}

func ToContainerStatsResponse(s db.ContainerStats) ContainerStatsResponse {
	return ContainerStatsResponse{
		ID:            s.ID,
		Name:          s.Name,
		CPUPercent:    s.CPUPercent,
		MemoryUsage:   s.MemoryUsage,
		MemoryLimit:   s.MemoryLimit,
		MemoryPercent: s.MemoryPercent,
		NetworkRx:     s.NetworkRx,
		NetworkTx:     s.NetworkTx,
	}
}

func ToContainerStatsListResponse(stats []db.ContainerStats) []ContainerStatsResponse {
	result := make([]ContainerStatsResponse, 0, len(stats))
	for _, s := range stats {
		result = append(result, ToContainerStatsResponse(s))
	}
	return result
}
