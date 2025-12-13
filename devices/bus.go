package devices

type BusImpl[T any] struct {
	Devices []T
}

func (b *BusImpl[T]) AddDevice(device T) {
	b.Devices = append(b.Devices, device)
}
