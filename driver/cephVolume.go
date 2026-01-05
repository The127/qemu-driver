package driver

type cephVolume struct {
	Id   string
	Pool string
	Name string
}

func (cephVolume) __volumeMarker() {}

func NewCephVolume(id string, pool string, name string) Volume {
	if len(id) > 24 {
		panic("id is too long")
	}
	return cephVolume{
		Id:   id,
		Pool: pool,
		Name: name,
	}
}
