package httphelper

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/flynn-examples/go-flynn-example/Godeps/_workspace/src/github.com/flynn/flynn/pkg/ctxhelper"
	log "github.com/flynn-examples/go-flynn-example/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
)

func NewRequestLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rw := w.(*ResponseWriter)

		reqID, _ := ctxhelper.RequestIDFromContext(rw.Context())
		componentName, _ := ctxhelper.ComponentNameFromContext(rw.Context())
		logger := log.New(log.Ctx{"component": componentName, "req_id": reqID})
		rw.ctx = ctxhelper.NewContextLogger(rw.Context(), logger)

		start := time.Now()
		var clientIP string
		clientIPs := strings.Split(req.Header.Get("X-Forwarded-For"), ",")
		if len(clientIPs) > 0 {
			clientIP = strings.TrimSpace(clientIPs[len(clientIPs)-1])
		}
		var err error
		if clientIP == "" {
			clientIP, _, err = net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				Error(w, err)
				return
			}
		}

		logger.Info("request started", "method", req.Method, "path", req.URL.Path, "client_ip", clientIP)

		handler.ServeHTTP(rw, req)

		logger.Info("request completed", "status", rw.Status(), "duration", time.Since(start))
	})
}
