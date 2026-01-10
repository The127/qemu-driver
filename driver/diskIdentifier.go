package driver

type DiskIdentifier struct {
	Vendor string
	Model  string
	Serial string
}

func DiskId(serial string) DiskIdentifier {
	return DiskIdentifier{
		Serial: serial,
	}
}

func DiskIdFull(vendor string, model string, serial string) DiskIdentifier {
	return DiskIdentifier{
		Vendor: vendor,
		Model:  model,
		Serial: serial,
	}
}
