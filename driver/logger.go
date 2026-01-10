package driver

type Logger interface {
	Logf(format string, v ...interface{})
}
