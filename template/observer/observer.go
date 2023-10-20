package observer

import (
	"log/slog"
	"net/http"
	"time"
)

func NewMiddleware(next http.Handler) Observer {
	return Observer{
		Next: next,
	}
}

type Observer struct {
	Next http.Handler
}

func (o Observer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l := &logger{ResponseWriter: w}
	start := time.Now()
	o.Next.ServeHTTP(l, r)

	slog.Info("",
		slog.Int("status", l.status),
		slog.Duration("duration", time.Since(start)),
		slog.String("method", r.Method),
		slog.Any("endpoint", r.URL),
	)
}

type logger struct {
	http.ResponseWriter
	status      int
	written     uint64
	wroteHeader bool
}

func (o *logger) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += uint64(n)

	return
}

func (o *logger) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}
