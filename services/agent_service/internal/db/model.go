package db

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
