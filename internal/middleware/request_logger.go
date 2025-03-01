package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/KengoWada/meetup-clone/internal/logger"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/go-chi/chi/v5/middleware"
)

var log = logger.Get()

func LoggerMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		start := time.Now()
		defer func() {
			end := time.Now()
			reqID := middleware.GetReqID(r.Context())

			if rvr := recover(); rvr != nil {
				log.Error().
					Str("type", "error").
					Timestamp().
					Interface("recover_info", rvr).
					Bytes("debug_stack", debug.Stack()).
					Msg("log system error")

				response.ErrorResponseInternalServerErr(ww, r, nil)
			}

			log.Info().
				Str("type", "access").
				Fields(map[string]any{
					"remote_ip":  r.RemoteAddr,
					"request_id": reqID,
					"url":        r.URL.Path,
					"proto":      r.Proto,
					"method":     r.Method,
					"user_agent": r.Header.Get("User-Agent"),
					"status":     ww.Status(),
					"latency_ms": float64(end.Sub(start).Nanoseconds()) / 1000000.0,
					"bytes_in":   r.Header.Get("Content-Length"),
					"bytes_out":  ww.BytesWritten(),
				}).
				Msg("Incoming Request")
		}()

		next.ServeHTTP(ww, r)
	}

	return http.HandlerFunc(fn)
}
