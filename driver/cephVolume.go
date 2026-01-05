package driver

type cephVolume struct {
	Serial string
	Pool   string
	Name   string
}

func (cephVolume) __volumeMarker() {}

func NewCephVolume(serial string, pool string, name string) Volume {
	return cephVolume{
		Serial: serial,
		Pool:   pool,
		Name:   name,
	}
}
