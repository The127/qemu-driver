package qmp

type PciDevice struct {
	Bus       int               `json:"bus"`
	Slot      int               `json:"slot"`
	Function  int               `json:"function"`
	ClassInfo PciDeviceClass    `json:"class_info"`
	Id        PciDeviceId       `json:"id"`
	Irq       int               `json:"irq"`
	IrqPin    int               `json:"irq_pin"`
	QdevId    string            `json:"qdev_id"`
	PciBridge PciBridgeInfo     `json:"pci_bridge"`
	Regions   []PciMemoryRegion `json:"regions"`
}

type PciDeviceClass struct {
	Class       int    `json:"class"`
	Description string `json:"desc"`
}

type PciDeviceId struct {
	Device          int `json:"device"`
	Vendor          int `json:"vendor"`
	Subsystem       int `json:"subsystem"`
	SubsystemVendor int `json:"subsystem-vendor"`
}

type PciBridgeInfo struct {
	Bus     PciBusInfo  `json:"bus"`
	Devices []PciDevice `json:"devices"`
}

type PciBusInfo struct {
	Number            int            `json:"number"`
	Secondary         int            `json:"secondary"`
	Subordinate       int            `json:"subordinate"`
	IoRange           PciMemoryRange `json:"io_range"`
	MemoryRange       PciMemoryRange `json:"memory_range"`
	PrefetchableRange PciMemoryRange `json:"prefetchable_range"`
}

type PciMemoryRange struct {
	Base  int `json:"base"`
	Limit int `json:"limit"`
}

type PciMemoryRegion struct {
	Bar       int    `json:"bar"`
	Type      string `json:"type"`
	Address   int    `json:"address"`
	Size      int    `json:"size"`
	Prefetch  bool   `json:"prefetch"`
	MemType64 bool   `json:"mem_type_64"`
}

type PciBus struct {
	Bus     int         `json:"bus"`
	Devices []PciDevice `json:"devices"`
}
