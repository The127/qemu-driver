package serial

type Bus interface {
	AddDevice(device BusDevice)
}
