package util

func BoolToOnOff(value bool) string {
	if value {
		return "on"
	} else {
		return "off"
	}
}
