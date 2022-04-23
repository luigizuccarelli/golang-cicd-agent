package service

import (
	"os"
	"testing"

	"github.com/microlib/logger/pkg/multi"

	"lmzsoftware.com/lzuccarelli/golang-cicd-engine/pkg/connectors"
	"lmzsoftware.com/lzuccarelli/golang-cicd-engine/pkg/schema"
)

func TestPipeline(t *testing.T) {

	// set up test env
	logger := multi.NewLogger(multi.HTML, multi.INFO)
	os.Setenv("CICD_CONSOLE_DIR", "../../tests/console")
	os.Setenv("WORKDIR", "../../tests/workdir")

	t.Run("ExecutePipeline : should pass", func(t *testing.T) {
		con := connectors.NewTestConnectors(logger, "normal")
		cicd := &schema.MapBinding{RepoName: "golang-test", RepoUrl: "test", Env: "dev"}
		err := ExecutePipeline(cicd, con)
		if err != nil {
			t.Fatalf("Should not fail : found error %v", err)
		}
	})

	t.Run("ExecutePipeline : should fail", func(t *testing.T) {
		con := connectors.NewTestConnectors(logger, "error")
		cicd := &schema.MapBinding{RepoName: "golang-test", RepoUrl: "test", Env: "dev"}
		err := ExecutePipeline(cicd, con)
		if err == nil {
			t.Fatalf("Should fail : did not find error ")
		}
	})

	t.Run("ExecutePipeline : should fail (rm forced error)", func(t *testing.T) {
		os.Setenv("WORKDIR", "./tmp")
		con := connectors.NewTestConnectors(logger, "error-rm")
		cicd := &schema.MapBinding{RepoName: "golang-test", RepoUrl: "test", Env: "dev"}
		err := ExecutePipeline(cicd, con)
		if err == nil {
			t.Fatalf("Should fail : did not find error ")
		}
	})

}
