package connectors

import (
	"bytes"
	"os/exec"

	"github.com/microlib/logger/pkg/multi"
)

// Connections struct - all backend connections in a common object
type Connectors struct {
	Logger *multi.Logger
}

func NewClientConnections(logger *multi.Logger) Clients {
	return &Connectors{Logger: logger}
}

func (c *Connectors) Error(msg string, val ...interface{}) {
	c.Logger.Errorf(msg, val...)
}

func (c *Connectors) Info(msg string, val ...interface{}) {
	c.Logger.Infof(msg, val...)
}

func (c *Connectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debugf(msg, val...)
}

func (c *Connectors) Trace(msg string, val ...interface{}) {
	c.Logger.Tracef(msg, val...)
}

func (c *Connectors) ExecOS(path string, command string, params []string, trim bool) (string, error) {
	var stdout, stderr bytes.Buffer
	var out string
	cmd := exec.Command(command, params...)
	cmd.Dir = path
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr := stdout.String()
	if err != nil {
		return stderr.String(), err
	}
	if trim {
		if len(outStr) > 1 {
			out = outStr[:len(outStr)-1]
		} else {
			out = outStr
		}
	} else {
		out = outStr
	}
	return out, nil
}
