package qmp

type MemorySummary struct {
	Base       uint64 `json:"base-memory"`
	Hotplugged uint64 `json:"plugged-memory"`
}

type MemoryInfo struct {
	Id           string `json:"id"`
	Addr         uint64 `json:"addr"`
	Size         uint64 `json:"size"`
	Slot         int    `json:"slot"`
	Node         int    `json:"node"`
	Memdev       string `json:"memdev"`
	Hotplugged   bool   `json:"hotplugged"`
	Hotpluggable bool   `json:"hotpluggable"`
}

type MemoryDevice = Wrapper[MemoryInfo]
