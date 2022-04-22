package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/microlib/logger/pkg/multi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"lmzsoftware.com/lzuccarelli/golang-cicd-engine/pkg/connectors"
	"lmzsoftware.com/lzuccarelli/golang-cicd-engine/pkg/handlers"
	"lmzsoftware.com/lzuccarelli/golang-cicd-engine/pkg/validator"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

var (
	logger       *multi.Logger
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)

// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
		// use this for cors
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept-Language")
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}

// startHttpServer - private function - http start
func startHttpServer(con connectors.Clients) *http.Server {
	srv := &http.Server{Addr: ":" + os.Getenv("SERVER_PORT")}

	r := mux.NewRouter()

	r.Use(prometheusMiddleware)
	r.Path("/metrics").Handler(promhttp.Handler())

	r.HandleFunc("/api/v1/webhook", func(w http.ResponseWriter, req *http.Request) {
		handlers.CICDHandler(w, req, con)
	}).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/v1/sys/info/isalive", handlers.IsAlive).Methods("GET")

	http.Handle("/", r)

	if err := srv.ListenAndServe(); err != nil {
		con.Error("Httpserver: ListenAndServe() error: " + err.Error())
	}

	return srv
}

func main() {
	if os.Getenv("LOG_LEVEL") == "" {
		logger = multi.NewLogger(multi.HTML, multi.TRACE)
	} else {
		logger = multi.NewLogger(multi.HTML, os.Getenv("LOG_LEVEL"))
	}
	err := validator.ValidateEnvars(logger)
	if err != nil {
		os.Exit(-1)
	}
	conn := connectors.NewClientConnections(logger)
	srv := startHttpServer(conn)
	logger.Info("Starting server on port " + srv.Addr)
}
