package driver

type StartOptions struct {
	CpuCount        int
	MemorySize      uint64
	DiskSize        uint64
	CloudInit       CloudInit
	Volumes         []Volume
	NetworkAdapters []NetworkAdapter
	VsockCid        uint32
}
