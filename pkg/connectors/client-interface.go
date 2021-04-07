package connectors

// Client Interface - used as a receiver and can be overriden for testing
type Clients interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
	ExecOS(path string, command string, params []string, trim bool) (string, error)
}
