package qmp

type CpuInfo struct {
	Index int `json:"cpu-index"`
}

type HotpluggableCpu struct {
	Type       string               `json:"type"`
	VCpusCount int                  `json:"vcpus-count"`
	QomPath    string               `json:"qom-path"`
	Props      HotpluggableCpuProps `json:"props"`
}

type HotpluggableCpuProps struct {
	SocketId int `json:"socket-id"`
	CoreId   int `json:"core-id"`
	ThreadId int `json:"thread-id"`
}
