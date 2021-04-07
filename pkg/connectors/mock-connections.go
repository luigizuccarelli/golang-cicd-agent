package connectors

import (
	"errors"
	"fmt"

	"github.com/microlib/logger/pkg/multi"
)

// Mock all connections
type MockConnectors struct {
	Logger *multi.Logger
	Flag   string
}

func (c *MockConnectors) Error(msg string, val ...interface{}) {
	c.Logger.Error(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) Info(msg string, val ...interface{}) {
	c.Logger.Info(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) Trace(msg string, val ...interface{}) {
	c.Logger.Trace(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) ExecOS(path string, command string, params []string, trim bool) (string, error) {
	// force error
	if c.Flag == "error-rm" && command == "rm" {
		return "", errors.New("forced rm error")
	}
	if c.Flag == "error" && command != "rm" {
		return "", errors.New("forced error")
	}
	return "testing", nil
}

func NewTestConnectors(logger *multi.Logger, flag string) Clients {
	conns := &MockConnectors{Logger: logger, Flag: flag}
	return conns
}
