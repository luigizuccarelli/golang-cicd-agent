package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/microlib/logger/pkg/multi"

	"lmzsoftware.com/lzuccarelli/golang-cicd-engine/pkg/connectors"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Inject (force) readAll test error")
}

func TestCICDHandler(t *testing.T) {

	logger := multi.NewLogger(multi.COLOR, multi.INFO)
	os.Setenv("WEBHOOK_SECRET", "Threefld2020")
	os.Setenv("REPO_MAPPING", "abc=xyz\ngolang-simple-oc4service=infra-golang\ntest=test")
	os.Setenv("CICD_CONSOLE_DIR", "../../tests/console")
	os.Setenv("WORKDIR", "../../tests/workdir")

	t.Run("IsAlive : should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/sys/info/isalive", nil)
		connectors.NewTestConnectors(logger, "normal")
		handler := http.HandlerFunc(IsAlive)
		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "IsAlive", rr.Code, STATUS))
		}
	})

	t.Run("CICDHandler : should pass (dev)", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		requestPayload, _ := ioutil.ReadFile("../../tests/merge.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/webhook", bytes.NewBuffer([]byte(requestPayload)))
		con := connectors.NewTestConnectors(logger, "normal")
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			CICDHandler(w, req, con)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "CICDHandler", rr.Code, STATUS))
		}
	})

	t.Run("CICDHandler : should pass (uat)", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		requestPayload, _ := ioutil.ReadFile("../../tests/uat-release.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/webhook", bytes.NewBuffer([]byte(requestPayload)))
		con := connectors.NewTestConnectors(logger, "normal")
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			CICDHandler(w, req, con)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "CICDHandler", rr.Code, STATUS))
		}
	})

	t.Run("CICDHandler : should pass (prod)", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		requestPayload, _ := ioutil.ReadFile("../../tests/prod-release.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/webhook", bytes.NewBuffer([]byte(requestPayload)))
		con := connectors.NewTestConnectors(logger, "normal")
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			CICDHandler(w, req, con)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "CICDHandler", rr.Code, STATUS))
		}
	})

	t.Run("CICDHandler : should fail (body content error)", func(t *testing.T) {
		var STATUS int = 500

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/webhook", errReader(0))
		conn := connectors.NewTestConnectors(logger, "normal")
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			CICDHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "CICDHandler ", rr.Code, STATUS))
		}
	})

	t.Run("CICDHandler : should fail (json to golang struct error)", func(t *testing.T) {
		var STATUS int = 500

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/webhook", bytes.NewBuffer([]byte("{ bad json")))
		conn := connectors.NewTestConnectors(logger, "normal")
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			CICDHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "CICDHandler ", rr.Code, STATUS))
		}
	})
	/*
		t.Run("CICDHandler : should fail (repo mapping)", func(t *testing.T) {
			var STATUS int = 500

			os.Setenv("REPO_MAPPING", "")
			requestPayload, _ := ioutil.ReadFile("../../tests/prod-release.json")
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/webhrook", bytes.NewBuffer([]byte(requestPayload)))
			conn := connectors.NewTestConnectors(logger, "normal")
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				CICDHandler(w, r, conn)
			})
			handler.ServeHTTP(rr, req)
			body, e := ioutil.ReadAll(rr.Body)
			if e != nil {
				t.Fatalf("Should not fail : found error %v", e)
			}
			logger.Trace(fmt.Sprintf("Response %s", string(body)))
			// ignore errors here
			if rr.Code != STATUS {
				t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "CICDHandler ", rr.Code, STATUS))
			}
		})
	*/

	t.Run("CICDHandler : should fail (api secret)", func(t *testing.T) {

		var STATUS int = 500
		os.Setenv("WEBHOOK_SECRET", "")
		requestPayload, _ := ioutil.ReadFile("../../tests/prod-release.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/webhook", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(logger, "normal")
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			CICDHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "CICDHandler ", rr.Code, STATUS))
		}
	})

}
