package connectors

import (
	"testing"

	"github.com/microlib/logger/pkg/multi"
)

func TestMockConnections(t *testing.T) {

	logger := multi.NewLogger(multi.COLOR, multi.TRACE)
	t.Run("MockConnections : ExecOS should pass", func(t *testing.T) {
		con := NewTestConnectors(logger, "normal")
		out, err := con.ExecOS(".", "ls", []string{"-ls"}, true)
		if err != nil {
			t.Fatalf("Should not fail : found error %v", err)
		}
		// testing logging
		con.Info("result : %s", out)
		con.Error("result : %s", out)
		con.Trace("result : %s", out)
		con.Debug("result : %s", out)
	})

	t.Run("Connections : ExecOS should fail", func(t *testing.T) {
		con := NewTestConnectors(logger, "error")
		out, err := con.ExecOS(".", "brrrr", []string{"-ls"}, true)
		if err == nil {
			t.Fatalf("Should fail : did not find error")
		}
		con.Info("result : %s", out)
		con.Error("err : %v", err)
	})

}
