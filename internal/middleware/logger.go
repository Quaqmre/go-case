package middleware

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func LoggerMiddleWare(l log.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			o := &responseObserver{ResponseWriter: w}
			h.ServeHTTP(o, r)
			addr := r.RemoteAddr
			if i := strings.LastIndex(addr, ":"); i != -1 {
				addr = addr[:i]
			}
			level.Info(l).Log("addr", addr, "detail", fmt.Sprintf("%s %s %s", r.Method, r.URL, r.Proto), "status", o.status, "user-agent", r.UserAgent)

		})
	}
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}
