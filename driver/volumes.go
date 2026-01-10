package driver

type Volume struct {
	id      DiskIdentifier
	options any
}

type cephVolumeOpts struct {
	pool string
	name string
}

func NewCephVolume(id DiskIdentifier, pool string, name string) Volume {
	return Volume{
		id: id,
		options: cephVolumeOpts{
			pool: pool,
			name: name,
		},
	}
}
