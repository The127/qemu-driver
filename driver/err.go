package driver

type RestartRequiredErr struct {
	msg string
}

func (e *RestartRequiredErr) Error() string {
	const text = "VM restart required"
	if e.msg == "" {
		return text
	}

	return text + ": " + e.msg
}

func newRestartRequiredErr(msg string) *RestartRequiredErr {
	return &RestartRequiredErr{msg: msg}
}
