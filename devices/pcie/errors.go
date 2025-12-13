package pcie

type deviceError int

func (n deviceError) Error() string {
	//goland:noinspection GoDirectComparisonOfErrors
	switch n {
	case NoHotplugError:
		return "this device can not be hotplugged"
	default:
		return "unknown error"
	}
}

const (
	NoHotplugError deviceError = iota
)
